package ai_platform_service_orchestrations

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	AnalysisTask = "analysis"
	AggTask      = "aggregation"
)

type TaskToExecute struct {
	WfID string                                           `json:"wfID"`
	Ou   org_users.OrgUser                                `json:"ou"`
	Tc   TaskContext                                      `json:"taskContext"`
	Wft  artemis_orchestrations.WorkflowTemplateData      `json:"wft"`
	Sg   *hera_search.SearchResultGroup                   `json:"sg"`
	Wr   *artemis_orchestrations.AIWorkflowAnalysisResult `json:"wr"`
}

type TaskContext struct {
	TaskName    string `json:"taskName"`
	TaskType    string `json:"taskType"`
	Model       string `json:"model"`
	TaskID      int    `json:"taskID"`
	ChunkOffset int    `json:"chunkOffset"`
}

func (z *ZeusAiPlatformServiceWorkflows) JsonOutputTaskWorkflow(ctx workflow.Context, tte TaskToExecute) (*ChatCompletionQueryResponse, error) {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(tte.Ou.OrgID, tte.WfID, "ZeusAiPlatformServiceWorkflows", "JsonOutputTaskWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return nil, err
	}

	var fullTaskDef []artemis_orchestrations.AITaskLibrary
	selectTaskCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(selectTaskCtx, z.SelectTaskDefinition, tte.Ou, tte.Sg.SourceTaskID).Get(selectTaskCtx, &fullTaskDef)
	if err != nil {
		logger.Error("failed to run task", "Error", err)
		return nil, err
	}
	if len(fullTaskDef) == 0 {
		return nil, nil
	}
	var jdef []*artemis_orchestrations.JsonSchemaDefinition
	for _, taskDef := range fullTaskDef {
		jdef = append(jdef, taskDef.Schemas...)
	}
	fd := artemis_orchestrations.ConvertToFuncDef(tte.Tc.TaskName, jdef)
	jsonTaskCtx := workflow.WithActivityOptions(ctx, ao)
	maxAttempts := ao.RetryPolicy.MaximumAttempts
	var aiResp *ChatCompletionQueryResponse
	for attempt := 0; attempt < int(maxAttempts); attempt++ {
		jsonTaskCtx = workflow.WithActivityOptions(ctx, ao)
		params := hera_openai.OpenAIParams{
			Model:              tte.Tc.Model,
			FunctionDefinition: fd,
		}
		err = workflow.ExecuteActivity(jsonTaskCtx, z.CreateJsonOutputModelResponse, tte.Ou, params).Get(jsonTaskCtx, &aiResp)
		if err != nil {
			logger.Error("failed to run analysis json", "Error", err)
			continue // Retry the activity
		}

		analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, tte.Ou, aiResp).Get(analysisCompCtx, &aiResp.ResponseTaskID)
		if err == nil {
			return nil, err
		}
		var m any
		if len(aiResp.Response.Choices) > 0 && len(aiResp.Response.Choices[0].Message.ToolCalls) > 0 {
			m, err = UnmarshallOpenAiJsonInterfaceSlice(params.FunctionDefinition.Name, aiResp)
			log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterfaceSlice failed")
			err = nil
		} else {
			m, err = UnmarshallOpenAiJsonInterface(params.FunctionDefinition.Name, aiResp)
			log.Err(err).Interface("m", m).Msg("UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
			err = nil
		}
		jsd := artemis_orchestrations.ConvertToJsonSchema(params.FunctionDefinition)
		aiResp.JsonResponseResults = artemis_orchestrations.AssignMapValuesMultipleJsonSchemasSlice(jsd, m)
		wr := artemis_orchestrations.AIWorkflowAnalysisResult{
			OrchestrationsID:      oj.OrchestrationID,
			ResponseID:            aiResp.ResponseTaskID,
			SourceTaskID:          tte.Tc.TaskID,
			ChunkOffset:           tte.Tc.ChunkOffset,
			IterationCount:        attempt,
			RunningCycleNumber:    tte.Wr.RunningCycleNumber,
			SearchWindowUnixStart: tte.Wr.SearchWindowUnixStart,
			SearchWindowUnixEnd:   tte.Wr.SearchWindowUnixEnd,
		}
		recordTaskCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(recordTaskCtx, z.SaveTaskOutput, wr, aiResp.JsonResponseResults).Get(recordTaskCtx, nil)
		if err != nil {
			logger.Error("failed to save task output", "Error", err)
			return nil, err
		}
	}

	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return nil, err
	}
	return aiResp, nil
}

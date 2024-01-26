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
	Ec   artemis_orchestrations.EvalContext               `json:"ec"`
	Tc   TaskContext                                      `json:"taskContext"`
	Wft  artemis_orchestrations.WorkflowTemplateData      `json:"wft"`
	Sg   *hera_search.SearchResultGroup                   `json:"sg"`
	Wr   *artemis_orchestrations.AIWorkflowAnalysisResult `json:"wr"`
}

type TaskContext struct {
	TaskName       string `json:"taskName"`
	TaskType       string `json:"taskType"`
	ResponseFormat string `json:"responseFormat"`
	Model          string `json:"model"`
	TaskID         int    `json:"taskID"`
	EvalID         int    `json:"evalID,omitempty"`
	Schemas        []*artemis_orchestrations.JsonSchemaDefinition
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

	jsonTaskCtx := workflow.WithActivityOptions(ctx, ao)
	maxAttempts := ao.RetryPolicy.MaximumAttempts
	var aiResp *ChatCompletionQueryResponse
	for attempt := 0; attempt < int(maxAttempts); attempt++ {
		jsonTaskCtx = workflow.WithActivityOptions(ctx, ao)
		fd := artemis_orchestrations.ConvertToFuncDef(tte.Tc.Schemas)
		params := hera_openai.OpenAIParams{
			Model:              tte.Tc.Model,
			Prompt:             tte.Sg.GetPromptBody(),
			FunctionDefinition: fd,
		}
		jsd := tte.Tc.Schemas
		tte.Wr.IterationCount = attempt
		err = workflow.ExecuteActivity(jsonTaskCtx, z.CreateJsonOutputModelResponse, tte.Ou, params).Get(jsonTaskCtx, &aiResp)
		if err != nil {
			logger.Error("failed to run analysis json", "Error", err)
			return nil, err
		}
		analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, tte.Ou, aiResp).Get(analysisCompCtx, &aiResp.ResponseTaskID)
		if err != nil {
			logger.Error("failed to record completion response", "Error", err)
			return nil, err
		}
		wr := tte.Wr
		wr.SourceTaskID = tte.Tc.TaskID
		wr.IterationCount = attempt
		wr.ResponseID = aiResp.ResponseTaskID
		var m any
		var anyErr error
		if len(aiResp.Response.Choices) > 0 && len(aiResp.Response.Choices[0].Message.ToolCalls) > 0 {
			m, anyErr = UnmarshallOpenAiJsonInterfaceSlice(params.FunctionDefinition.Name, aiResp)
			// ok no err
			if anyErr != nil {
				log.Err(anyErr).Interface("m", m).Msg("1_UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterfaceSlice failed")
			}
		} else {
			m, anyErr = UnmarshallOpenAiJsonInterface(params.FunctionDefinition.Name, aiResp)
			if anyErr != nil {
				log.Err(anyErr).Interface("m", m).Msg("2_UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
			}
		}
		var tmpResp [][]*artemis_orchestrations.JsonSchemaDefinition
		if anyErr == nil {
			tmpResp, anyErr = artemis_orchestrations.AssignMapValuesMultipleJsonSchemasSlice(jsd, m)
		}
		if m == nil || len(tmpResp) == 0 {
			continue
		}
		if len(tmpResp) > 0 && len(tmpResp[0]) <= 0 {
			continue
		}
		if anyErr != nil {
			log.Err(anyErr).Interface("m", m).Msg("JsonOutputTaskWorkflow: AssignMapValuesMultipleJsonSchemasSlice: failed")
			tte.Wr.SkipAnalysis = true
			if tte.Tc.EvalID > 0 {
				evrr := artemis_orchestrations.AIWorkflowEvalResultResponse{
					EvalID:             tte.Tc.EvalID,
					WorkflowResultID:   wr.WorkflowResultID,
					ResponseID:         wr.ResponseID,
					EvalIterationCount: wr.IterationCount,
				}
				recordEvalResCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(recordEvalResCtx, z.SaveEvalResponseOutput, evrr).Get(recordEvalResCtx, nil)
				if err != nil {
					logger.Error("failed to save eval resp id", "Error", err)
					return nil, err
				}
			} else {
				recordTaskCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(recordTaskCtx, z.SaveTaskOutput, tte.Wr, aiResp.Response).Get(recordTaskCtx, &aiResp.WorkflowResultID)
				if err != nil {
					logger.Error("failed to save task output", "Error", err)
					return nil, err
				}
			}
			continue
		}
		aiResp.JsonResponseResults = append(aiResp.JsonResponseResults, tmpResp...)
		if tte.Tc.EvalID > 0 {
			evrr := artemis_orchestrations.AIWorkflowEvalResultResponse{
				EvalID:             tte.Tc.EvalID,
				WorkflowResultID:   wr.WorkflowResultID,
				ResponseID:         wr.ResponseID,
				EvalIterationCount: wr.IterationCount,
			}
			recordEvalResCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(recordEvalResCtx, z.SaveEvalResponseOutput, evrr).Get(recordEvalResCtx, nil)
			if err != nil {
				logger.Error("failed to save eval resp id", "Error", err)
				return nil, err
			}
		} else {
			recordTaskCtx := workflow.WithActivityOptions(ctx, ao)
			tte.Wr.SkipAnalysis = false
			err = workflow.ExecuteActivity(recordTaskCtx, z.SaveTaskOutput, tte.Wr, aiResp.JsonResponseResults).Get(recordTaskCtx, &aiResp.WorkflowResultID)
			if err != nil {
				logger.Error("failed to save task output", "Error", err)
				return nil, err
			}
		}
		break
	}

	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return nil, err
	}
	return aiResp, nil
}

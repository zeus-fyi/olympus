package ai_platform_service_orchestrations

import (
	"time"

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
	WfID     string                                      `json:"wfID"`
	Ou       org_users.OrgUser                           `json:"ou"`
	TaskType string                                      `json:"taskType"`
	Wft      artemis_orchestrations.WorkflowTemplateData `json:"wft"`
	Sg       *hera_search.SearchResultGroup              `json:"sg"`
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
	pr := &PromptReduction{
		PromptReductionSearchResults: &PromptReductionSearchResults{
			InSearchGroup: tte.Sg,
		},
	}
	taskName := ""
	switch tte.TaskType {
	case AnalysisTask:
		taskName = tte.Wft.AnalysisTaskName
		pr.Model = tte.Wft.AnalysisModel
		pr.TokenOverflowStrategy = tte.Wft.AnalysisTokenOverflowStrategy
		pr.PromptReductionText.InPromptBody = tte.Wft.AnalysisPrompt
	case AggTask:
		if tte.Wft.AggTaskName == nil || tte.Wft.AggModel == nil || tte.Wft.AggTokenOverflowStrategy == nil || tte.Wft.AggPrompt == nil {
			return nil, nil
		}
		taskName = *tte.Wft.AggTaskName
		pr.Model = *tte.Wft.AggModel
		pr.TokenOverflowStrategy = *tte.Wft.AggTokenOverflowStrategy
		pr.PromptReductionText.InPromptBody = *tte.Wft.AggPrompt
	default:
		return nil, nil
	}
	chunkedTaskCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(chunkedTaskCtx, z.TokenOverflowReduction, tte.Ou, pr).Get(chunkedTaskCtx, &pr)
	if err != nil {
		logger.Error("failed to run token overflow task", "Error", err)
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
	var aiResp *ChatCompletionQueryResponse
	fd := artemis_orchestrations.ConvertToFuncDef(taskName, jdef)
	jsonTaskCtx := workflow.WithActivityOptions(ctx, ao)
	params := hera_openai.OpenAIParams{
		Model:              pr.Model,
		FunctionDefinition: fd,
	}
	err = workflow.ExecuteActivity(jsonTaskCtx, z.CreateJsonOutputModelResponse, tte.Ou, params).Get(jsonTaskCtx, &aiResp)
	if err != nil {
		logger.Error("failed to run analysis json", "Error", err)
		return nil, err
	}
	analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, tte.Ou, aiResp).Get(analysisCompCtx, &aiResp.ResponseTaskID)
	if err != nil {
		logger.Error("failed to save analysis response", "Error", err)
		return nil, err
	}
	wr := artemis_orchestrations.AIWorkflowAnalysisResult{
		//OrchestrationsID:      oj.OrchestrationID,
		//ResponseID:            analysisRespId,
		//SourceTaskID:          analysisInst.AnalysisTaskID,
		//RunningCycleNumber:    i,
		//SearchWindowUnixStart: window.UnixStartTime,
		//SearchWindowUnixEnd:   window.UnixEndTime,
		IterationCount: 1,
	}
	recordAnalysisCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(recordAnalysisCtx, z.SaveTaskOutput, wr).Get(recordAnalysisCtx, nil)
	if err != nil {
		logger.Error("failed to save analysis", "Error", err)
		return nil, err
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return nil, err
	}
	return aiResp, nil
}

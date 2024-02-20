package ai_platform_service_orchestrations

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	AnalysisTask = "analysis"
	AggTask      = "aggregation"
)

type TaskToExecute struct {
	Ec  artemis_orchestrations.EvalContext               `json:"ec"`
	Wft artemis_orchestrations.WorkflowTemplateData      `json:"wft"`
	Wr  *artemis_orchestrations.AIWorkflowAnalysisResult `json:"wr"`
}

type TaskContext struct {
	TaskName                           string                                                       `json:"taskName"`
	TaskType                           string                                                       `json:"taskType"`
	Temperature                        float32                                                      `json:"temperature"`
	MarginBuffer                       float64                                                      `json:"marginBuffer"`
	Prompt                             string                                                       `json:"prompt,omitempty"`
	TokenOverflowStrategy              string                                                       `json:"tokenOverflowStrategy"`
	ResponseFormat                     string                                                       `json:"responseFormat"`
	Model                              string                                                       `json:"model"`
	EvalModel                          string                                                       `json:"evalModel"`
	WorkflowResultID                   int                                                          `json:"workflowResultID"`
	TaskID                             int                                                          `json:"taskID"`
	EvalID                             int                                                          `json:"evalID,omitempty"`
	Retrieval                          artemis_orchestrations.RetrievalItem                         `json:"retrieval,omitempty"`
	TriggerActionsApproval             artemis_orchestrations.TriggerActionsApproval                `json:"triggerActionsApproval,omitempty"`
	AIWorkflowTriggerResultApiResponse artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse `json:"aiWorkflowTriggerResultApiResponse,omitempty"`
	EvalSchemas                        []*artemis_orchestrations.JsonSchemaDefinition               `json:"evalSchemas,omitempty"`
	Schemas                            []*artemis_orchestrations.JsonSchemaDefinition               `json:"schemas,omitempty"`
	JsonResponseResults                []artemis_orchestrations.JsonSchemaDefinition                `json:"jsonResponseResults,omitempty"`
	RetryPolicy                        *temporal.RetryPolicy                                        `json:"retryPolicy"`
	WorkflowTemplateData               artemis_orchestrations.WorkflowTemplateData                  `json:"workflowTemplateData"`
}

func (z *ZeusAiPlatformServiceWorkflows) JsonOutputTaskWorkflow(ctx workflow.Context, mb *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	var canSkip bool
	if mb.Tc.JsonResponseResults != nil && len(mb.Tc.JsonResponseResults) > 0 && len(mb.Tc.EvalSchemas) > 0 {
		jrs := mb.Tc.JsonResponseResults
		canSkip = mb.Tc.EvalModel == mb.Tc.Model
		for _, sv := range mb.Tc.EvalSchemas {
			if !CheckSchemaIDsAndValidFields(sv.SchemaID, jrs) {
				canSkip = false
				break
			}
		}
	}
	if canSkip {
		return mb, nil
	}
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    25,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(mb.Ou.OrgID, mb.Wsr.ChildWfID, "ZeusAiPlatformServiceWorkflows", "JsonOutputTaskWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return nil, err
	}
	jsonTaskCtx := workflow.WithActivityOptions(ctx, ao)
	maxAttempts := ao.RetryPolicy.MaximumAttempts
	var aiResp *ChatCompletionQueryResponse
	var feedback error
	for attempt := 0; attempt < int(maxAttempts); attempt++ {
		log.Info().Int("attempt", attempt).Msg("JsonOutputTaskWorkflow: attempt")
		jsonTaskCtx = workflow.WithActivityOptions(ctx, ao)
		feedbackPrompt := ""
		if feedback != nil {
			feedbackPrompt = fmt.Sprintf("Please fix your answer or make best assumptions on data structure to fix this error: %s. This is attempt number: %d", feedback.Error(), attempt)
			feedback = nil
		}
		params := hera_openai.OpenAIParams{
			Model:              mb.Tc.Model,
			FunctionDefinition: artemis_orchestrations.ConvertToFuncDef(mb.Tc.Schemas),
			Temperature:        mb.Tc.Temperature,
			SystemPromptExt:    feedbackPrompt,
		}
		wfa := artemis_orchestrations.AIWorkflowAnalysisResult{
			OrchestrationID:       mb.Oj.OrchestrationID,
			SourceTaskID:          mb.Tc.TaskID,
			IterationCount:        attempt,
			ChunkOffset:           mb.Wsr.ChunkOffset,
			RunningCycleNumber:    mb.Wsr.RunCycle,
			SearchWindowUnixStart: mb.Window.UnixStartTime,
			SearchWindowUnixEnd:   mb.Window.UnixEndTime,
		}
		mb.Wsr.IterationCount = attempt
		err = workflow.ExecuteActivity(jsonTaskCtx, z.CreateJsonOutputModelResponse, mb, params).Get(jsonTaskCtx, &aiResp)
		if err != nil {
			logger.Error("failed to run analysis json", "Error", err)
			feedback = err
			continue
		}
		wfa.ResponseID = aiResp.ResponseID
		aiResp.WorkflowResultID = wfa.WorkflowResultID
		log.Info().Int("attempt", attempt).Interface("len(aiResp.JsonResponseResults)", len(aiResp.JsonResponseResults)).Msg("JsonOutputTaskWorkflow: done")
		if mb.Tc.EvalID > 0 {
			evrr := artemis_orchestrations.AIWorkflowEvalResultResponse{
				EvalID:             mb.Tc.EvalID,
				WorkflowResultID:   wfa.WorkflowResultID,
				ResponseID:         wfa.ResponseID,
				EvalIterationCount: wfa.IterationCount,
			}
			recordEvalResCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(recordEvalResCtx, z.SaveEvalResponseOutput, evrr).Get(recordEvalResCtx, &aiResp.EvalResultID)
			if err != nil {
				logger.Error("failed to save eval resp id", "Error", err)
				return nil, err
			}
		} else {
			recordTaskCtx := workflow.WithActivityOptions(ctx, ao)
			wfa.SkipAnalysis = false
			err = workflow.ExecuteActivity(recordTaskCtx, z.SaveTaskOutput, wfa, aiResp.JsonResponseResults).Get(recordTaskCtx, &aiResp.WorkflowResultID)
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
	mb.Tc.JsonResponseResults = aiResp.JsonResponseResults
	mb.Tc.WorkflowResultID = aiResp.WorkflowResultID
	mb.Tc.Schemas = aiResp.Schemas
	return mb, nil
}

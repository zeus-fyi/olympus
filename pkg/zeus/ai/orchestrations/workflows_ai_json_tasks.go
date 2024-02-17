package ai_platform_service_orchestrations

import (
	"fmt"
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
	WfID        string                                           `json:"wfID"`
	Ou          org_users.OrgUser                                `json:"ou"`
	Ec          artemis_orchestrations.EvalContext               `json:"ec"`
	Tc          TaskContext                                      `json:"taskContext"`
	Wft         artemis_orchestrations.WorkflowTemplateData      `json:"wft"`
	Sg          *hera_search.SearchResultGroup                   `json:"sg"`
	Wr          *artemis_orchestrations.AIWorkflowAnalysisResult `json:"wr"`
	RetryPolicy *temporal.RetryPolicy                            `json:"retryPolicy"`
}

type TaskContext struct {
	TaskName                           string                                        `json:"taskName"`
	TaskType                           string                                        `json:"taskType"`
	Temperature                        float32                                       `json:"temperature"`
	ResponseFormat                     string                                        `json:"responseFormat"`
	Model                              string                                        `json:"model"`
	TaskID                             int                                           `json:"taskID"`
	EvalID                             int                                           `json:"evalID,omitempty"`
	Retrieval                          artemis_orchestrations.RetrievalItem          `json:"retrieval,omitempty"`
	TriggerActionsApproval             artemis_orchestrations.TriggerActionsApproval `json:"triggerActionsApproval,omitempty"`
	AIWorkflowTriggerResultApiResponse artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse
	Schemas                            []*artemis_orchestrations.JsonSchemaDefinition
}

func (z *ZeusAiPlatformServiceWorkflows) JsonOutputTaskWorkflow(ctx workflow.Context, tte TaskToExecute) (*ChatCompletionQueryResponse, error) {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 10,
			MaximumAttempts:    25,
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
	var feedback error

	for attempt := 0; attempt < int(maxAttempts); attempt++ {
		log.Info().Int("attempt", attempt).Msg("JsonOutputTaskWorkflow: attempt")
		jsonTaskCtx = workflow.WithActivityOptions(ctx, ao)
		fd := artemis_orchestrations.ConvertToFuncDef(tte.Tc.Schemas)
		feedbackPrompt := ""
		if feedback != nil {
			feedbackPrompt = fmt.Sprintf("Please fix your answer or make best assumptions on data structure to fix this error: %s. This is attempt number: %d", feedback.Error(), attempt)
			feedback = nil
		}
		params := hera_openai.OpenAIParams{
			Model:              tte.Tc.Model,
			Prompt:             tte.Sg.GetPromptBody(),
			FunctionDefinition: fd,
			Temperature:        tte.Tc.Temperature,
			SystemPromptExt:    feedbackPrompt,
		}
		jsd := tte.Tc.Schemas
		tte.Wr.IterationCount = attempt
		err = workflow.ExecuteActivity(jsonTaskCtx, z.CreateJsonOutputModelResponse, tte.Ou, params).Get(jsonTaskCtx, &aiResp)
		if err != nil {
			logger.Error("failed to run analysis json", "Error", err)
			return nil, err
		}
		analysisCompCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(analysisCompCtx, z.RecordCompletionResponse, tte.Ou, aiResp).Get(analysisCompCtx, &aiResp.ResponseID)
		if err != nil {
			logger.Error("failed to record completion response", "Error", err)
			return nil, err
		}
		wr := tte.Wr
		wr.SourceTaskID = tte.Tc.TaskID
		wr.IterationCount = attempt
		wr.ResponseID = aiResp.ResponseID
		var m any
		var anyErr error
		if len(aiResp.Response.Choices) > 0 && len(aiResp.Response.Choices[0].Message.ToolCalls) > 0 {
			m, anyErr = UnmarshallOpenAiJsonInterfaceSlice(params.FunctionDefinition.Name, aiResp)
			// ok no err
			if anyErr != nil {
				log.Err(anyErr).Interface("m", m).Interface("resp", aiResp.Response.Choices).Msg("1_UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterfaceSlice failed")
				logger.Error("1_UnmarshallFilteredMsgIdsFromAiJson", "Error", err, "m", m, "resp", aiResp.Response.Choices)
			}
		} else {
			m, anyErr = UnmarshallOpenAiJsonInterface(params.FunctionDefinition.Name, aiResp)
			if anyErr != nil {
				log.Err(anyErr).Interface("m", m).Interface("resp", aiResp.Response.Choices).Msg("2_UnmarshallFilteredMsgIdsFromAiJson: UnmarshallOpenAiJsonInterface failed")
				logger.Error("2_UnmarshallFilteredMsgIdsFromAiJson", "Error", err, "m", m, "resp", aiResp.Response.Choices)
			}
		}
		var tmpResp []artemis_orchestrations.JsonSchemaDefinition
		if anyErr == nil {
			tmpResp, anyErr = artemis_orchestrations.AssignMapValuesMultipleJsonSchemasSlice(jsd, m)
			if anyErr != nil {
				log.Err(anyErr).Interface("m", m).Interface("jsd", jsd).Msg("AssignMapValuesMultipleJsonSchemasSlice: UnmarshallOpenAiJsonInterface failed")
				feedback = anyErr
			}
		} else {
			feedback = anyErr
		}
		if anyErr == nil {
			aiResp.JsonResponseResults = tmpResp
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
				tte.Ec.AIWorkflowEvalResultResponse = evrr
				recordEvalResCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(recordEvalResCtx, z.SaveEvalResponseOutput, evrr).Get(recordEvalResCtx, &aiResp.EvalResultID)
				if err != nil {
					logger.Error("failed to save eval resp id", "Error", err)
					return nil, err
				}
				continue
			} else {
				afv := InputDataAnalysisToAgg{
					ChatCompletionQueryResponse: aiResp,
					SearchResultGroup:           tte.Sg,
				}
				recordTaskCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(recordTaskCtx, z.SaveTaskOutput, tte.Wr, afv).Get(recordTaskCtx, &aiResp.WorkflowResultID)
				if err != nil {
					logger.Error("failed to save task output", "Error", err)
					return nil, err
				}
			}
			continue
		}
		aiResp.JsonResponseResults = tmpResp
		log.Info().Int("attempt", attempt).Interface("len(aiResp.JsonResponseResults)", len(aiResp.JsonResponseResults)).Msg("JsonOutputTaskWorkflow: done")
		if tte.Tc.EvalID > 0 {
			evrr := artemis_orchestrations.AIWorkflowEvalResultResponse{
				EvalID:             tte.Tc.EvalID,
				WorkflowResultID:   wr.WorkflowResultID,
				ResponseID:         wr.ResponseID,
				EvalIterationCount: wr.IterationCount,
			}
			recordEvalResCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(recordEvalResCtx, z.SaveEvalResponseOutput, evrr).Get(recordEvalResCtx, &aiResp.EvalResultID)
			if err != nil {
				logger.Error("failed to save eval resp id", "Error", err)
				return nil, err
			}
		} else {
			recordTaskCtx := workflow.WithActivityOptions(ctx, ao)
			tte.Wr.SkipAnalysis = false
			afv := InputDataAnalysisToAgg{
				ChatCompletionQueryResponse: aiResp,
				SearchResultGroup:           tte.Sg,
			}
			err = workflow.ExecuteActivity(recordTaskCtx, z.SaveTaskOutput, tte.Wr, afv).Get(recordTaskCtx, &aiResp.WorkflowResultID)
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

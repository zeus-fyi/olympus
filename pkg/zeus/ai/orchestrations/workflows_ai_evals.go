package ai_platform_service_orchestrations

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type MbChildSubProcessParams struct {
	WfID                     string                                          `json:"wfID"`
	Ou                       org_users.OrgUser                               `json:"ou"`
	WfExecParams             artemis_orchestrations.WorkflowExecParams       `json:"wfExecParams"`
	Oj                       artemis_orchestrations.OrchestrationJob         `json:"oj"`
	RunCycle                 int                                             `json:"runCycle"`
	Window                   artemis_orchestrations.Window                   `json:"window"`
	WorkflowResult           artemis_orchestrations.AIWorkflowAnalysisResult `json:"workflowResult"`
	AnalysisEvalActionParams *EvalActionParams                               `json:"analysisEvalActionParams,omitempty"`
}

type EvalActionParams struct {
	WorkflowTemplateData artemis_orchestrations.WorkflowTemplateData `json:"parentProcess"`
	ParentOutputToEval   *ChatCompletionQueryResponse                `json:"parentOutputToEval"`
	EvalFns              []artemis_orchestrations.EvalFnDB           `json:"evalFns"`
	SearchResultGroup    *hera_search.SearchResultGroup              `json:"searchResultsGroup,omitempty"`
}

func (z *ZeusAiPlatformServiceWorkflows) RunAiWorkflowAutoEvalProcess(ctx workflow.Context, mb *MbChildSubProcessParams, cpe *EvalActionParams) error {
	if cpe == nil || mb == nil {
		return nil
	}
	logger := workflow.GetLogger(ctx)
	aoAiAct := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    5,
		},
	}

	for _, evalFn := range cpe.EvalFns {
		evalFnMetricsLookupCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		var evalFnMetrics []artemis_orchestrations.EvalFn
		err := workflow.ExecuteActivity(evalFnMetricsLookupCtx, z.EvalLookup, mb.Ou, evalFn.EvalID).Get(evalFnMetricsLookupCtx, &evalFnMetrics)
		if err != nil {
			logger.Error("failed to get eval info", "Error", err)
			return err
		}
		for _, evalFnWithMetrics := range evalFnMetrics {
			fd, ferr := TransformEvalMetricsToJSONSchema(evalFnWithMetrics.EvalMetrics)
			if ferr != nil {
				logger.Error("failed to transform eval metrics to json schema", "Error", ferr)
				return ferr
			}
			evalParams := hera_openai.OpenAIParams{
				Model: evalFn.EvalModel,
				FunctionDefinition: openai.FunctionDefinition{
					Name:        evalFnWithMetrics.EvalName,
					Description: evalFnWithMetrics.EvalName,
					Parameters:  fd,
				},
			}
			var emr *artemis_orchestrations.EvalMetricsResults
			evCtx := artemis_orchestrations.EvalContext{
				EvalID:                evalFn.EvalID,
				OrchestrationID:       mb.Oj.OrchestrationID,
				SourceTaskID:          cpe.ParentOutputToEval.ResponseTaskID,
				RunningCycleNumber:    mb.RunCycle,
				SearchWindowUnixStart: mb.Window.UnixStartTime,
				SearchWindowUnixEnd:   mb.Window.UnixEndTime,
				WorkflowResultID:      mb.WorkflowResult.WorkflowResultID,
			}

			switch strings.ToLower(evalFnWithMetrics.EvalType) {
			case "model":
				var cr *ChatCompletionQueryResponse
				modelScoredJsonCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				err = workflow.ExecuteActivity(modelScoredJsonCtx, z.CreateJsonOutputModelResponse, mb.Ou, evalParams).Get(modelScoredJsonCtx, &cr)
				if err != nil {
					logger.Error("failed to get eval info", "Error", err)
					return err
				}
				if cr == nil || len(cr.Response.Choices) == 0 {
					continue
				}

				var evalCompletionID int
				evalJsonCompCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				err = workflow.ExecuteActivity(evalJsonCompCtx, z.RecordCompletionResponse, mb.Ou, cr).Get(evalJsonCompCtx, &evalCompletionID)
				if err != nil {
					logger.Error("failed to save eval json response", "Error", err)
					return err
				}
				recordEvalResponseCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				evrr := artemis_orchestrations.AIWorkflowEvalResultResponse{
					WorkflowResultID: mb.WorkflowResult.WorkflowResultID,
					EvalID:           evalFn.EvalID,
					ResponseID:       evalCompletionID,
				}
				err = workflow.ExecuteActivity(recordEvalResponseCtx, z.SaveEvalResponseOutput, evrr).Get(recordEvalResponseCtx, nil)
				if err != nil {
					logger.Error("failed to save eval response relationship", "Error", err)
					return err
				}
				m := make(map[string]interface{})
				for _, cho := range cr.Response.Choices {
					for _, tvr := range cho.Message.ToolCalls {
						if tvr.Function.Name == evalFnWithMetrics.EvalName {
							err = json.Unmarshal([]byte(tvr.Function.Arguments), &m)
							if err != nil {
								logger.Error("failed to unmarshal json", "Error", err)
								return err
							}
						}
					}
				}
				if len(m) == 0 {
					logger.Warn("failed to get eval info", "Response", m)
					continue
				}
				evalModelScoredJsonCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				err = workflow.ExecuteActivity(evalModelScoredJsonCtx, z.EvalModelScoredJsonOutput, m, &evalFnWithMetrics).Get(evalModelScoredJsonCtx, &emr)
				if err != nil {
					logger.Error("failed to get score eval", "Error", err)
					return err
				}
			case "api":
				// TODO, complete this, should attach a retrieval option? use that for the scoring?
				//retrievalCtx := workflow.WithActivityOptions(ctx, ao)
				//err = workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, ou, analysisInst, window).Get(retrievalCtx, &sr)
				//if err != nil {
				//	logger.Error("failed to run retrieval", "Error", err)
				//	return err
				//}
				//var routes []iris_models.RouteInfo
				//retrievalWebCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				//err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, ou, analysisInst).Get(retrievalWebCtx, &routes)
				//if err != nil {
				//	logger.Error("failed to run retrieval", "Error", err)
				//	return err
				//}
				//for _, route := range routes {
				//	fetchedResult := &hera_search.SearchResult{}
				//	retrievalWebTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				//	err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.AiWebRetrievalTask, ou, analysisInst, route).Get(retrievalWebTaskCtx, &fetchedResult)
				//	if err != nil {
				//		logger.Error("failed to run retrieval", "Error", err)
				//		return err
				//	}
				//	if fetchedResult != nil && len(fetchedResult.Value) > 0 {
				//		sr = append(sr, *fetchedResult)
				//	}
				//}
				//cr := &ChatCompletionQueryResponse{}
				//apiScoredJsonCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				//err = workflow.ExecuteActivity(apiScoredJsonCtx, z.SendResponseToApiForScoresInJson, mb.Ou, evalParams).Get(apiScoredJsonCtx, &cr)
				//if err != nil {
				//	logger.Error("failed to get eval info", "Error", err)
				//	return err
				//}
			}
			if emr == nil {
				continue
			}
			for _, er := range emr.EvalMetricsResults {
				// in the eval stage, if any filter fails, skip the analysis.
				if er.EvalState == "filter" && ((er.EvalMetricResult == "pass" && er.EvalResultOutcome == false) || (er.EvalMetricResult == "fail" && er.EvalResultOutcome == true)) {
					mb.WorkflowResult.SkipAnalysis = true
					recordAnalysisCtx := workflow.WithActivityOptions(ctx, aoAiAct)
					err = workflow.ExecuteActivity(recordAnalysisCtx, z.SaveTaskOutput, mb.WorkflowResult).Get(recordAnalysisCtx, nil)
					if err != nil {
						logger.Error("failed to save analysis skip", "Error", err)
						return err
					}
					break
				}
			}
			saveEvalResultsCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			emr.EvalContext = evCtx
			err = workflow.ExecuteActivity(saveEvalResultsCtx, z.SaveEvalMetricResults, emr).Get(saveEvalResultsCtx, nil)
			if err != nil {
				logger.Error("failed to save eval metric results", "Error", err)
				return err
			}
			suffix := strings.Split(uuid.New().String(), "-")[0]
			childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
				WorkflowID:               mb.Oj.OrchestrationName + "-eval-trigger-" + strconv.Itoa(mb.RunCycle) + suffix,
				WorkflowExecutionTimeout: mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
			}
			childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
			tar := TriggerActionsWorkflowParams{
				Emr: emr,
				Mb:  mb,
			}
			err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RunTriggerActions, tar).Get(childAnalysisCtx, nil)
			if err != nil {
				logger.Error("failed to execute child run trigger actions workflow", "Error", err)
				return err
			}
		}
	}
	return nil
}

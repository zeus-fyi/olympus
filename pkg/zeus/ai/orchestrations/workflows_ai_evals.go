package ai_platform_service_orchestrations

import (
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
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
	TaskToExecute        TaskToExecute                               `json:"tte,omitempty"`
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
			MaximumAttempts:    10,
		},
	}

	evalsFnsMap := make(map[int]*artemis_orchestrations.EvalFn)
	var evalFnsAgg []artemis_orchestrations.EvalFn
	for ei, _ := range cpe.EvalFns {
		if cpe.EvalFns[ei].EvalID == 0 {
			continue
		}
		var efs []artemis_orchestrations.EvalFn
		evalFnMetricsLookupCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err := workflow.ExecuteActivity(evalFnMetricsLookupCtx, z.EvalLookup, mb.Ou, cpe.EvalFns[ei].EvalID).Get(evalFnMetricsLookupCtx, &efs)
		if err != nil {
			logger.Error("failed to get eval info", "Error", err)
			return err
		}
		for _, ef := range efs {
			evalsFnsMap[aws.IntValue(ef.EvalID)] = &ef
		}
		evalFnsAgg = append(evalFnsAgg, efs...)
	}
	for evFnIndex, _ := range evalFnsAgg {
		if evalFnsAgg[evFnIndex].EvalID == nil {
			continue
		}
		var emr *artemis_orchestrations.EvalMetricsResults
		evCtx := artemis_orchestrations.EvalContext{
			EvalID:                   aws.IntValue(evalFnsAgg[evFnIndex].EvalID),
			AIWorkflowAnalysisResult: mb.WorkflowResult,
			EvalIterationCount:       0,
		}
		cpe.TaskToExecute.Ec = evCtx
		switch strings.ToLower(evalFnsAgg[evFnIndex].EvalType) {
		case "model":
			if cpe.ParentOutputToEval != nil && cpe.ParentOutputToEval.JsonResponseResults != nil &&
				cpe.ParentOutputToEval.Params.Model == aws.StringValue(evalFnsAgg[evFnIndex].EvalModel) &&
				copyMatchingJsonResponsesFieldValuesFromResp(cpe.ParentOutputToEval.JsonResponseResults, evalFnsAgg[evFnIndex].SchemasMap) {
			} else {
				wfID := mb.Oj.OrchestrationName + "-automated-model-scored-evals-" + strconv.Itoa(mb.RunCycle)
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					WorkflowID:               wfID,
					WorkflowExecutionTimeout: mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
				}
				if len(evalFnsAgg[evFnIndex].Schemas) == 0 {
					continue
				}
				cpe.TaskToExecute.Tc.Schemas = evalFnsAgg[evFnIndex].Schemas
				cpe.TaskToExecute.Tc.Model = aws.StringValue(evalFnsAgg[evFnIndex].EvalModel)
				cpe.ParentOutputToEval = &ChatCompletionQueryResponse{}
				childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
				err := workflow.ExecuteChildWorkflow(childAnalysisCtx, z.JsonOutputTaskWorkflow, cpe.TaskToExecute).Get(childAnalysisCtx, &cpe.ParentOutputToEval)
				if err != nil {
					logger.Error("failed to execute analysis json workflow", "Error", err)
					return err
				}
			}
			evalModelScoredJsonCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			err := workflow.ExecuteActivity(evalModelScoredJsonCtx, z.EvalModelScoredJsonOutput, cpe.ParentOutputToEval.JsonResponseResults, &evalFnsAgg[evFnIndex]).Get(evalModelScoredJsonCtx, &emr)
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
		emr.EvalContext = evCtx
		saveEvalResultsCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err := workflow.ExecuteActivity(saveEvalResultsCtx, z.SaveEvalMetricResults, emr).Get(saveEvalResultsCtx, nil)
		if err != nil {
			logger.Error("failed to save eval metric results", "Error", err)
			return err
		}
		suffix := strings.Split(uuid.New().String(), "-")[0]
		wfID := mb.Oj.OrchestrationName + "-eval-trigger-" + strconv.Itoa(mb.RunCycle) + suffix
		childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
			WorkflowID:               wfID,
			WorkflowExecutionTimeout: mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
		}
		childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
		tar := TriggerActionsWorkflowParams{
			Emr: emr,
			Mb:  mb,
		}
		err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.CreateTriggerActionsWorkflow, tar).Get(childAnalysisCtx, nil)
		if err != nil {
			logger.Error("failed to execute child run trigger actions workflow", "Error", err)
			return err
		}
	}
	return nil
}

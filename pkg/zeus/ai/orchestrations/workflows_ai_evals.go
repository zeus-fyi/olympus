package ai_platform_service_orchestrations

import (
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

const (
	evalModelScoredJsonOutput = "model"
	evalModelScoredViaApi     = "api"
)

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
			MaximumAttempts:    25,
		},
	}
	evalsFnsMap := make(map[int]*artemis_orchestrations.EvalFn)
	var evalFnsAgg []artemis_orchestrations.EvalFn
	for ei, _ := range cpe.EvalFns {
		if cpe.EvalFns[ei].EvalID == 0 {
			continue
		}
		log.Info().Int("evalID", cpe.EvalFns[ei].EvalID).Msg("evalID")
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
		evCtx := artemis_orchestrations.EvalContext{EvalID: aws.IntValue(evalFnsAgg[evFnIndex].EvalID), AIWorkflowAnalysisResult: mb.WorkflowResult, EvalIterationCount: 0}
		cpe.TaskToExecute.Ec = evCtx
		switch strings.ToLower(evalFnsAgg[evFnIndex].EvalType) {
		case evalModelScoredJsonOutput:
			wfID := mb.Oj.OrchestrationName + "-automated-model-scored-evals-" + strconv.Itoa(mb.RunCycle)
			childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{WorkflowID: wfID, WorkflowExecutionTimeout: mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize, RetryPolicy: aoAiAct.RetryPolicy}
			if len(evalFnsAgg[evFnIndex].Schemas) == 0 {
				continue
			}
			var canSkip bool
			if cpe.ParentOutputToEval != nil && cpe.ParentOutputToEval.JsonResponseResults != nil && len(cpe.ParentOutputToEval.JsonResponseResults) > 0 && len(evalFnsAgg[evFnIndex].Schemas) > 0 {
				jrs := cpe.ParentOutputToEval.JsonResponseResults
				evs := evalFnsAgg[evFnIndex].Schemas
				evmModelName := aws.StringValue(evalFnsAgg[evFnIndex].EvalModel)
				canSkip = evmModelName == cpe.ParentOutputToEval.Params.Model
				for _, sv := range evs {
					if !CheckSchemaIDsAndValidFields(sv.SchemaID, jrs) {
						canSkip = false
						break
					}
				}
			}
			cpe.TaskToExecute.Tc.EvalID = aws.IntValue(evalFnsAgg[evFnIndex].EvalID)
			cpe.TaskToExecute.Tc.Schemas = evalFnsAgg[evFnIndex].Schemas
			cpe.TaskToExecute.Tc.Model = aws.StringValue(evalFnsAgg[evFnIndex].EvalModel)
			if !canSkip {
				childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
				err := workflow.ExecuteChildWorkflow(childAnalysisCtx, z.JsonOutputTaskWorkflow, cpe.TaskToExecute).Get(childAnalysisCtx, &cpe.ParentOutputToEval)
				if err != nil {
					logger.Error("failed to execute analysis json workflow", "Error", err)
					return err
				}
			}
			if cpe.ParentOutputToEval.JsonResponseResults == nil {
				logger.Warn("json response results are nil, skipping eval", "evalID", cpe.TaskToExecute)
				continue
			}
			efin := evalFnsAgg[evFnIndex]
			jrevr := cpe.ParentOutputToEval.JsonResponseResults
			evalModelScoredJsonCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			err := workflow.ExecuteActivity(evalModelScoredJsonCtx, z.EvalModelScoredJsonOutput, efin, jrevr).Get(evalModelScoredJsonCtx, &emr)
			if err != nil {
				logger.Error("failed to get score eval", "Error", err)
				return err
			}
			evCtx.JsonResponseResults = cpe.ParentOutputToEval.JsonResponseResults
			evCtx.EvaluatedJsonResponses = emr.EvaluatedJsonResponses
		case evalModelScoredViaApi:
		}
		if emr == nil {
			log.Warn().Msg("emr is nil, skipping save eval metric results")
			continue
		}
		emr.EvalContext = evCtx
		saveEvalResultsCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err := workflow.ExecuteActivity(saveEvalResultsCtx, z.SaveEvalMetricResults, emr).Get(saveEvalResultsCtx, nil)
		if err != nil {
			logger.Error("failed to save eval metric results", "Error", err)
			return err
		}
		cpe.TaskToExecute.Ec = evCtx
		tar := TriggerActionsWorkflowParams{Emr: emr, Mb: mb, Cpe: cpe}
		childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
			WorkflowID:               mb.Oj.OrchestrationName + "-eval-trigger-" + strconv.Itoa(mb.RunCycle) + strings.Split(uuid.New().String(), "-")[0],
			WorkflowExecutionTimeout: mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
			RetryPolicy:              aoAiAct.RetryPolicy,
		}
		childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
		err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.CreateTriggerActionsWorkflow, tar).Get(childAnalysisCtx, nil)
		if err != nil {
			logger.Error("failed to execute child run trigger actions workflow", "Error", err)
			return err
		}
	}
	return nil
}

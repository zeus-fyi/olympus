package ai_platform_service_orchestrations

import (
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type MbChildSubProcessParams struct {
	WfID         string                                        `json:"wfID"`
	Ou           org_users.OrgUser                             `json:"ou"`
	WfExecParams artemis_orchestrations.WorkflowExecParams     `json:"wfExecParams"`
	Oj           artemis_orchestrations.OrchestrationJob       `json:"oj"`
	Window       artemis_orchestrations.Window                 `json:"window"`
	Wsr          artemis_orchestrations.WorkflowStageReference `json:"wsr"`
	Tc           TaskContext                                   `json:"taskContext"`
}

const (
	evalModelScoredJsonOutput = "model"
	evalModelScoredViaApi     = "api"
)

type RunAiWorkflowAutoEvalProcessInputs struct {
	Mb *MbChildSubProcessParams `json:"mb,omitempty"`
}

func (z *ZeusAiPlatformServiceWorkflows) RunAiWorkflowAutoEvalProcess(ctx workflow.Context, mb *MbChildSubProcessParams) error {
	if mb == nil || mb.Tc.EvalID == 0 {
		return nil
	}
	logger := workflow.GetLogger(ctx)
	aoAiAct := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    25,
		},
	}
	evalsFnsMap := make(map[int]*artemis_orchestrations.EvalFn)
	var evalFnsAgg []artemis_orchestrations.EvalFn
	log.Info().Int("evalID", mb.Tc.EvalID).Interface("taskType", mb.Tc.TaskType).Msg("RunAiWorkflowAutoEvalProcess: evalID")
	var efs []artemis_orchestrations.EvalFn
	evalFnMetricsLookupCtx := workflow.WithActivityOptions(ctx, aoAiAct)
	err := workflow.ExecuteActivity(evalFnMetricsLookupCtx, z.EvalLookup, mb.Ou, mb.Tc.EvalID).Get(evalFnMetricsLookupCtx, &efs)
	if err != nil {
		logger.Error("failed to get eval info", "Error", err)
		return err
	}
	for _, ef := range efs {
		if ef.EvalID == nil {
			continue
		}
		evalsFnsMap[aws.IntValue(ef.EvalID)] = &ef
	}
	evalFnsAgg = append(evalFnsAgg, efs...)
	log.Info().Interface("RunAiWorkflowAutoEvalProcess: len(evalFnsAgg)", len(evalFnsAgg)).Msg("evalFnsAgg")
	for evFnIndex, _ := range evalFnsAgg {
		switch strings.ToLower(evalFnsAgg[evFnIndex].EvalType) {
		case evalModelScoredJsonOutput:
			if len(evalFnsAgg[evFnIndex].Schemas) == 0 {
				continue
			}
			mb.Tc.EvalID = aws.IntValue(evalFnsAgg[evFnIndex].EvalID)
			mb.Tc.EvalSchemas = evalFnsAgg[evFnIndex].Schemas
			mb.Tc.EvalModel = aws.StringValue(evalFnsAgg[evFnIndex].EvalModel)
			mb.Wsr.ChildWfID = mb.Oj.OrchestrationName + "-automated-model-scored-evals-" + strconv.Itoa(mb.Wsr.RunCycle) + "-ind-" + strconv.Itoa(evFnIndex)
			childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{WorkflowID: mb.Wsr.ChildWfID, WorkflowExecutionTimeout: mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize, RetryPolicy: aoAiAct.RetryPolicy}
			childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
			err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.JsonOutputTaskWorkflow, mb).Get(childAnalysisCtx, &mb)
			if err != nil {
				logger.Error("failed to execute analysis json workflow", "Error", err)
				return err
			}
			if mb.Tc.JsonResponseResults == nil {
				logger.Warn("json response results are nil, skipping eval", "evalID", mb.Tc.EvalID)
				continue
			}
			evalModelScoredJsonCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			err = workflow.ExecuteActivity(evalModelScoredJsonCtx, z.EvalModelScoredJsonOutput, evalFnsAgg[evFnIndex], mb).Get(evalModelScoredJsonCtx, &mb)
			if err != nil {
				logger.Error("failed to get score eval", "Error", err)
				return err
			}
		case evalModelScoredViaApi:
		}
		childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
			WorkflowID:               mb.Oj.OrchestrationName + "-eval-trigger-" + strconv.Itoa(mb.Wsr.RunCycle) + "-ind-" + strconv.Itoa(evFnIndex),
			WorkflowExecutionTimeout: mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
			RetryPolicy:              aoAiAct.RetryPolicy,
		}
		childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
		err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.CreateTriggerActionsWorkflow, mb).Get(childAnalysisCtx, nil)
		if err != nil {
			logger.Error("failed to execute child run trigger actions workflow", "Error", err)
			return err
		}
	}
	return nil
}

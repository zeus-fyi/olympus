package kronos_helix

import (
	"strings"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

func (k *KronosWorkflow) Mockingbird(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
	}
	aCtx := workflow.WithActivityOptions(ctx, ao)
	var ojs []artemis_orchestrations.OrchestrationJob
	err := workflow.ExecuteActivity(aCtx, k.SelectOrchestrationsByGroupNameAndType, mockingbird, cronjob).Get(aCtx, &ojs)
	if err != nil {
		logger.Error("failed to get internal assignments", "Error", err)
		return err
	}
	for _, oj := range ojs {
		var inst Instructions
		instCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(instCtx, k.GetInstructionsFromJob, oj).Get(instCtx, &inst)
		if err != nil {
			logger.Error("failed to get alert assignment from instructions", "Error", err)
			return err
		}

		switch strings.ToLower(oj.Type) {
		case alerts:
		case monitoring:
		case cronjob:
			childWorkflowOptions := workflow.ChildWorkflowOptions{
				TaskQueue:         KronosHelixTaskQueue,
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			}
			childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
			childWfFuture := workflow.ExecuteChildWorkflow(childCtx, "KronosCronJob", inst, CalculatePollCycles(kronosLoopInterval, inst.CronJob.PollInterval))
			var childWE workflow.Execution
			if err = childWfFuture.GetChildWorkflowExecution().Get(childCtx, &childWE); err != nil {
				logger.Error("Failed to get child workflow execution", "Error", err)
				return err
			}
		}
	}
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		TaskQueue:         KronosHelixTaskQueue,
		ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
	}
	childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
	childWfFuture := workflow.ExecuteChildWorkflow(childCtx, "Yang", kronosLoopInterval)
	var childWE workflow.Execution
	if err = childWfFuture.GetChildWorkflowExecution().Get(childCtx, &childWE); err != nil {
		logger.Error("Failed to get child workflow execution", "Error", err)
		return err
	}
	return nil
}

package kronos_helix

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (k *KronosWorkflow) Monitor(ctx workflow.Context, oj *artemis_orchestrations.OrchestrationJob, mi Instructions, pollCycles int) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    mi.Monitors.PollInterval,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
		},
	}
	failureCount := 0
	for i := 0; i < pollCycles; i++ {
		healthCtx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(healthCtx, k.CheckEndpointHealth, mi, failureCount).Get(healthCtx, &failureCount)
		if err != nil {
			logger.Error("failed to execute triggered alert", "Error", err)
			// You can decide if you want to return the error or continue monitoring.
			return err
		}
		if failureCount >= mi.Monitors.AlertFailureThreshold {
			alertCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(alertCtx, k.ExecuteTriggeredAlert, CreateHealthMonitorAlertEvent(mi.Monitors.ServiceName)).Get(alertCtx, nil)
			if err != nil {
				logger.Error("failed to execute triggered alert", "Error", err)
				// You can decide if you want to return the error or continue monitoring.
				return err
			}
			childWorkflowOptions := workflow.ChildWorkflowOptions{
				TaskQueue:         KronosHelixTaskQueue,
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
				RetryPolicy:       ao.RetryPolicy,
			}
			childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
			childWfFuture := workflow.ExecuteChildWorkflow(childCtx, "OrchestrationChildProcessReset", &oj, mi)
			var childWE workflow.Execution
			if err = childWfFuture.GetChildWorkflowExecution().Get(childCtx, &childWE); err != nil {
				logger.Error("Failed to get child workflow execution", "Error", err)
				return err
			}
			return nil
		}
		err = workflow.Sleep(ctx, mi.Monitors.PollInterval)
		if err != nil {
			logger.Error("failed to sleep", "Error", err)
			return err
		}
	}
	return nil
}

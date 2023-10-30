package kronos_helix

import (
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type KronosWorkflow struct {
	temporal_base.Workflow
	KronosActivities
}

const (
	kronosLoopInterval = 10 * time.Minute
	alerts             = "alerts"
	monitoring         = "monitoring"
)

func NewKronosWorkflow() KronosWorkflow {
	deployWf := KronosWorkflow{
		Workflow:         temporal_base.Workflow{},
		KronosActivities: KronosActivities{},
	}
	return deployWf
}

func CalculatePollCycles(intervalResetTime time.Duration, pollInterval time.Duration) int {
	return int(intervalResetTime / pollInterval)
}

func (k *KronosWorkflow) GetWorkflows() []interface{} {
	return []interface{}{
		k.Yin, k.Yang, k.SignalFlow,
		k.OrchestrationChildProcessReset,
		k.Monitor, k.CronJob,
	}
}

// SignalFlow should be used to place new control flows on the helix
func (k *KronosWorkflow) SignalFlow(ctx workflow.Context) error {
	return nil
}

// Yin should send commands and execute actions
func (k *KronosWorkflow) Yin(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
	}
	aCtx := workflow.WithActivityOptions(ctx, ao)
	var ojs []artemis_orchestrations.OrchestrationJob
	err := workflow.ExecuteActivity(aCtx, k.GetInternalAssignments).Get(aCtx, &ojs)
	if err != nil {
		logger.Error("failed to get internal assignments", "Error", err)
		return err
	}
	for _, oj := range ojs {
		var pdV2Event *pagerduty.V2Event
		var inst Instructions
		instCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(instCtx, k.GetInstructionsFromJob, oj).Get(instCtx, &inst)
		if err != nil {
			logger.Error("failed to get alert assignment from instructions", "Error", err)
			return err
		}

		switch oj.Type {
		case alerts:
			alertAssignmentCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(alertAssignmentCtx, k.GetAlertAssignmentFromInstructions, inst).Get(alertAssignmentCtx, &pdV2Event)
			if err != nil {
				logger.Error("failed to get alert assignment from instructions", "Error", err)
				return err
			}
			if pdV2Event != nil && pdV2Event.DedupKey != "" {
				alertCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(alertCtx, k.ExecuteTriggeredAlert, pdV2Event).Get(alertCtx, nil)
				if err != nil {
					logger.Error("failed to execute triggered alert", "Error", err)
					return err
				}
				childWorkflowOptions := workflow.ChildWorkflowOptions{
					TaskQueue:         KronosHelixTaskQueue,
					ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
				}
				childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
				childWfFuture := workflow.ExecuteChildWorkflow(childCtx, "OrchestrationChildProcessReset", &oj, inst)
				var childWE workflow.Execution
				if err = childWfFuture.GetChildWorkflowExecution().Get(childCtx, &childWE); err != nil {
					logger.Error("Failed to get child workflow execution", "Error", err)
					return err
				}
			}
		case monitoring:
			childWorkflowOptions := workflow.ChildWorkflowOptions{
				TaskQueue:         KronosHelixTaskQueue,
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			}
			childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
			childWfFuture := workflow.ExecuteChildWorkflow(childCtx, "Monitor", &oj, inst, CalculatePollCycles(kronosLoopInterval, inst.Monitors.PollInterval))
			var childWE workflow.Execution
			if err = childWfFuture.GetChildWorkflowExecution().Get(childCtx, &childWE); err != nil {
				logger.Error("Failed to get child workflow execution", "Error", err)
				return err
			}
		case Cronjob:
			logger.Info("Cronjob", "inst", inst)
			childWorkflowOptions := workflow.ChildWorkflowOptions{
				TaskQueue:         KronosHelixTaskQueue,
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
			}
			childCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
			childWfFuture := workflow.ExecuteChildWorkflow(childCtx, "CronJob", inst, CalculatePollCycles(kronosLoopInterval, inst.CronJob.PollInterval))
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

func (k *KronosWorkflow) OrchestrationChildProcessReset(ctx workflow.Context, oj *artemis_orchestrations.OrchestrationJob, inst Instructions) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour * 730, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 10,
			BackoffCoefficient: 2,
			MaximumInterval:    time.Minute * 5,
		},
	}
	aCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(aCtx, k.UpdateAndMarkOrchestrationInactive, &oj).Get(aCtx, nil)
	if err != nil {
		logger.Error("failed to execute triggered alert", "Error", err)
		return err
	}
	err = workflow.Sleep(aCtx, inst.Trigger.ResetAlertAfterTimeDuration)
	if err != nil {
		logger.Error("failed to sleep", "Error", err)
		return err
	}
	err = workflow.ExecuteActivity(aCtx, k.UpdateAndMarkOrchestrationActive, &oj).Get(aCtx, nil)
	if err != nil {
		logger.Error("failed to execute triggered alert", "Error", err)
		return err
	}
	return nil
}

// Yang should check, status, & react to changes
func (k *KronosWorkflow) Yang(ctx workflow.Context, waitTime time.Duration) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Minute,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
		},
	}
	err := workflow.Sleep(ctx, waitTime)
	if err != nil {
		logger.Error("failed to sleep", "Error", err)
		return err
	}
	aCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(aCtx, k.Recycle).Get(aCtx, nil)
	if err != nil {
		logger.Error("failed to recycle", "Error", err)
		return err
	}
	return nil
}

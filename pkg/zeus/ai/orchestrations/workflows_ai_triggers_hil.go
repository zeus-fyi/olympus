package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

/*
1. should update state to approved or should update state to rejected, then end workflow
2. should lookup trigger information
	- social media engagement (twitter, reddit, discord, telegram)
		execute API calls

*/

func (z *ZeusAiPlatformServiceWorkflows) RunApprovedTriggerActions(ctx workflow.Context, tar TriggerActionsWorkflowParams) error {
	if tar.Emr == nil || tar.Mb == nil {
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

	// if conditions are met, create or update the trigger action
	recordTriggerCondCtx := workflow.WithActivityOptions(ctx, aoAiAct)
	var ta *artemis_orchestrations.TriggerAction
	err := workflow.ExecuteActivity(recordTriggerCondCtx, z.CreateOrUpdateTriggerActionToExec, tar.Mb, ta).Get(recordTriggerCondCtx, nil)
	if err != nil {
		logger.Error("failed to create or update trigger action", "Error", err)
		return err
	}
	switch ta.TriggerEnv {
	case "social-media-engagement":
		childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
			//WorkflowID:               mb.Oj.OrchestrationName + "-eval-trigger-" + strconv.Itoa(mb.RunCycle) + suffix,
			//WorkflowExecutionTimeout: mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
		}
		childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
		err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RunApprovedSocialMediaTriggerActionsWorkflow, tar).Get(childAnalysisCtx, nil)
		if err != nil {
			logger.Error("failed to execute child run trigger actions workflow", "Error", err)
			return err
		}
	}
	return nil
}

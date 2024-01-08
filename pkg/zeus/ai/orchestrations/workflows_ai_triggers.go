package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type TriggerActionsWorkflowParams struct {
	Emr *artemis_orchestrations.EvalMetricsResults
	Mb  *MbChildSubProcessParams
}

func (z *ZeusAiPlatformServiceWorkflows) RunTriggerActions(ctx workflow.Context, tar TriggerActionsWorkflowParams) error {
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

	triggerEvalsLookupCtx := workflow.WithActivityOptions(ctx, aoAiAct)
	var triggerActions []artemis_orchestrations.TriggerAction
	err := workflow.ExecuteActivity(triggerEvalsLookupCtx, z.LookupEvalTriggerConditions, tar.Mb.Ou, tar.Emr.EvalContext.EvalID).Get(triggerEvalsLookupCtx, &triggerActions)
	if err != nil {
		logger.Error("failed to get eval info", "Error", err)
		return err
	}

	for _, triggerAction := range triggerActions {
		var ta *artemis_orchestrations.TriggerAction
		err = workflow.ExecuteActivity(ctx, z.CheckEvalTriggerCondition, &triggerAction, tar.Emr).Get(ctx, &ta)
		if err != nil {
			logger.Error("failed to check eval trigger condition", "Error", err)
			return err
		}
		if ta == nil {
			continue
		}
		err = workflow.ExecuteActivity(ctx, z.CreateOrUpdateTriggerActionToExec, tar.Mb, ta).Get(ctx, nil)
		if err != nil {
			logger.Error("failed to create or update trigger action", "Error", err)
			return err
		}
	}
	return nil
}

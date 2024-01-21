package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (z *ZeusAiPlatformServiceWorkflows) RunApprovedSocialMediaTriggerActionsWorkflow(ctx workflow.Context, tar TriggerActionsWorkflowParams) error {
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
	smType := ""
	switch smType {
	case twitterPlatform:
		socialMediaExecCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(socialMediaExecCtx, z.SocialTweetTask, tar.Mb.Ou, tar.Emr.EvalContext.EvalID).Get(socialMediaExecCtx, nil)
		if err != nil {
			logger.Error("failed to exec twitter api call", "Error", err)
			return err
		}
	case redditPlatform:
		socialMediaExecCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(socialMediaExecCtx, z.SocialRedditTask, tar.Mb.Ou, tar.Emr.EvalContext.EvalID).Get(socialMediaExecCtx, nil)
		if err != nil {
			logger.Error("failed to exec reddit api call", "Error", err)
			return err
		}
	case discordPlatform:
		socialMediaExecCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(socialMediaExecCtx, z.SocialDiscordTask, tar.Mb.Ou, tar.Emr.EvalContext.EvalID).Get(socialMediaExecCtx, nil)
		if err != nil {
			logger.Error("failed to exec discord api call", "Error", err)
			return err
		}
	case telegramPlatform:
		socialMediaExecCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(socialMediaExecCtx, z.SocialTelegramTask, tar.Mb.Ou, tar.Emr.EvalContext.EvalID).Get(socialMediaExecCtx, nil)
		if err != nil {
			logger.Error("failed to exec telegram api call", "Error", err)
			return err
		}
	}
	return nil
}

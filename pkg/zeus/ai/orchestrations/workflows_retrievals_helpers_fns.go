package ai_platform_service_orchestrations

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func getPlatform(cp *MbChildSubProcessParams) string {
	platform := cp.Tc.Retrieval.RetrievalPlatform
	if cp.Tc.TriggerActionsApproval.TriggerAction == apiApproval {
		platform = apiApproval
	}
	if cp.Tc.EvalID <= 0 || cp.Tc.TriggerActionsApproval.ApprovalID <= 0 {
		switch platform {
		case twitterPlatform, redditPlatform, discordPlatform, telegramPlatform:
		default:
			platform = webPlatform
		}
	}
	return platform
}

func getRetActRetryPolicy(cp *MbChildSubProcessParams) workflow.ActivityOptions {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 24, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.5,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    100,
		},
	}
	if cp.Tc.RetryPolicy != nil {
		ao.RetryPolicy = cp.Tc.RetryPolicy
	}
	return ao
}

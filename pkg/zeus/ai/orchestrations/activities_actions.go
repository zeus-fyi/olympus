package ai_platform_service_orchestrations

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (z *ZeusAiPlatformActivities) LookupEvalTriggerConditions(ctx context.Context) error {
	return nil
}

// SendTriggerActionRequestForApproval sends the action request to the user for human in-the-loop approval
func (z *ZeusAiPlatformActivities) SendTriggerActionRequestForApproval(ctx context.Context) error {
	return nil
}

func (z *ZeusAiPlatformActivities) CreateOrUpdateTriggerActionRequestForApproval(ctx context.Context) error {
	return nil
}

func (z *ZeusAiPlatformActivities) CheckEvalTriggerCondition(ctx context.Context, act artemis_orchestrations.TriggerAction) error {
	return nil
}

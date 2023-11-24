package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hera_discord "github.com/zeus-fyi/olympus/pkg/hera/discord"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (h *ZeusAiPlatformServiceWorkflows) AiIngestDiscordWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, searchGroupName string, cm hera_discord.ChannelMessages) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "ZeusAiPlatformServiceWorkflows", "AiIngestDiscordWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}
	searchQueryCtx := workflow.WithActivityOptions(ctx, ao)
	var sq *hera_search.DiscordSearchResultWrapper
	err = workflow.ExecuteActivity(searchQueryCtx, h.SelectDiscordSearchQuery, ou, searchGroupName).Get(searchQueryCtx, &sq)
	if err != nil {
		logger.Error("failed to execute SelectDiscordSearchQuery", "Error", err)
		// You can decide if you want to return the error or continue monitoring.
		return err
	}
	if sq == nil || sq.SearchID == 0 {
		logger.Info("no new tweets found")
		return nil
	}
	insertMessagesCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(insertMessagesCtx, h.InsertIncomingDiscordDataFromSearch, sq.SearchID, cm).Get(insertMessagesCtx, nil)
	if err != nil {
		logger.Error("failed to execute InsertIncomingDiscordDataFromSearch", "Error", err)
		// You can decide if you want to return the error or continue monitoring.
		return err
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
	return nil
}

func (h *ZeusAiPlatformServiceWorkflows) AiFetchDataToIngestDiscordWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, searchGroupName string) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "ZeusAiPlatformServiceWorkflows", "AiFetchDataToIngestDiscordWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}

	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
	return nil
}

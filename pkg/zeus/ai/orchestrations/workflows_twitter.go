package ai_platform_service_orchestrations

import (
	"time"

	"github.com/cvcio/twitter"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (h *ZeusAiPlatformServiceWorkflows) AiIngestTwitterWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, groupName string) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}
	// todo allow user orgs ids
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "ZeusAiPlatformServiceWorkflows", "AiIngestTwitterWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}
	insertMsgCtx := workflow.WithActivityOptions(ctx, ao)
	var sq *hera_search.TwitterSearchQuery
	err = workflow.ExecuteActivity(insertMsgCtx, h.SelectTwitterSearchQuery, ou, groupName).Get(insertMsgCtx, &sq)
	if err != nil {
		logger.Error("failed to execute SelectTwitterSearchQuery", "Error", err)
		// You can decide if you want to return the error or continue monitoring.
		return err
	}
	var tweets []*twitter.Tweet
	searchCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(searchCtx, h.SearchTwitterUsingQuery, sq).Get(searchCtx, &tweets)
	if err != nil {
		logger.Error("failed to execute InsertEmailIfNew", "Error", err)
		// You can decide if you want to return the error or continue monitoring.
		return err
	}
	if len(tweets) == 0 {
		logger.Info("no new tweets found")
		return nil
	}
	insertTweetsCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(insertTweetsCtx, h.InsertIncomingTweetsFromSearch, sq.SearchID, tweets).Get(insertTweetsCtx, &tweets)
	if err != nil {
		logger.Error("failed to execute InsertIncomingTweetsFromSearch", "Error", err)
		// You can decide if you want to return the error or continue monitoring.
		return err
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for services", "Error", err)
		return err
	}
	return nil
}

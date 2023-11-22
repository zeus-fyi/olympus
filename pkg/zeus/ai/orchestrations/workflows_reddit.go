package ai_platform_service_orchestrations

import (
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (h *ZeusAiPlatformServiceWorkflows) AiIngestRedditWorkflow(ctx workflow.Context, wfID string, ou org_users.OrgUser, groupName string) error {
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
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "ZeusAiPlatformServiceWorkflows", "AiIngestRedditWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return err
	}
	// Activity to select the Reddit search query. Insert this part after the UpsertAssignmentActivity
	selectRedditQueryCtx := workflow.WithActivityOptions(ctx, ao)
	var redditSearchQuery *hera_search.RedditSearchQuery
	err = workflow.ExecuteActivity(selectRedditQueryCtx, h.SelectRedditSearchQuery, ou, groupName).Get(selectRedditQueryCtx, &redditSearchQuery)
	if err != nil {
		logger.Error("failed to select Reddit search query", "Error", err)
		return err
	}

	if redditSearchQuery == nil {
		logger.Info("no Reddit search query found")
		return nil
	}
	redditCtx := workflow.WithActivityOptions(ctx, ao)
	lpo := &reddit.ListOptions{Limit: redditSearchQuery.MaxResults, After: redditSearchQuery.FullPostId}
	var redditPosts []*reddit.Post
	err = workflow.ExecuteActivity(redditCtx, h.SearchRedditNewPostsUsingSubreddit, redditSearchQuery.Query, lpo).Get(redditCtx, &redditPosts)
	if err != nil {
		logger.Error("failed to fetch new Reddit posts", "Error", err)
		return err
	}
	// Add the InsertIncomingRedditDataFromSearch activity here
	if redditPosts == nil || len(redditPosts) == 0 {
		logger.Info("no new Reddit posts found")
		return nil
	}
	insertRedditDataCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(insertRedditDataCtx, h.InsertIncomingRedditDataFromSearch, redditSearchQuery.SearchID, redditPosts).Get(insertRedditDataCtx, nil)
	if err != nil {
		logger.Error("failed to insert incoming Reddit data from search", "Error", err)
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

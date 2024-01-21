package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type TaskToExecute struct {
	WfID string                                      `json:"wfID"`
	Ou   org_users.OrgUser                           `json:"ou"`
	Wft  artemis_orchestrations.WorkflowTemplateData `json:"wft"`
	Sg   *hera_search.SearchResultGroup              `json:"sg"`
}

func (z *ZeusAiPlatformServiceWorkflows) RetrievalsWorkflow(ctx workflow.Context, tte TaskToExecute) (*hera_search.SearchResultGroup, error) {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(tte.Ou.OrgID, tte.WfID, "ZeusAiPlatformServiceWorkflows", "RetrievalsWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return nil, err
	}
	switch tte.Wft.RetrievalPlatform {
	case twitterPlatform, redditPlatform, discordPlatform, telegramPlatform:
		retrievalCtx := workflow.WithActivityOptions(ctx, ao)
		var sr []hera_search.SearchResult
		err = workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, tte.Ou.OrgID, tte.Wft, tte.Sg.Window).Get(retrievalCtx, &sr)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
		tte.Sg.SearchResults = append(tte.Sg.SearchResults, sr...)
	case webPlatform:
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, tte.Ou.OrgID, tte.Wft).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
		for _, route := range routes {
			fetchedResult := &hera_search.SearchResult{}
			retrievalWebTaskCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.AiWebRetrievalTask, tte.Ou.OrgID, tte.Wft, route).Get(retrievalWebTaskCtx, &fetchedResult)
			if err != nil {
				logger.Error("failed to run retrieval", "Error", err)
				return nil, err
			}
			if fetchedResult != nil && len(fetchedResult.Value) > 0 {
				tte.Sg.SearchResults = append(tte.Sg.SearchResults, *fetchedResult)
			}
		}
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return nil, err
	}
	return tte.Sg, nil
}

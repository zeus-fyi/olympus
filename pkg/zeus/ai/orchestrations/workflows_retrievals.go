package ai_platform_service_orchestrations

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

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
		err = workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, tte.Ou, tte.Wft, tte.Sg.Window).Get(retrievalCtx, &tte.Sg)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
		tte.Sg.SourceTaskID = tte.Wft.AnalysisTaskID
	case apiApproval, webPlatform:
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, tte.Ou.OrgID, tte.Wft).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
		var echoResp []echo.Map
		if tte.Wft.RetrievalPlatform == apiApproval {
			payloadMaps := artemis_orchestrations.CreateMapInterfaceFromAssignedSchemaFields(tte.Ec.JsonResponseResults)
			for _, m := range payloadMaps {
				echoMap := echo.Map{}
				for k, v := range m {
					echoMap[k] = v
				}
				echoResp = append(echoResp, echoMap)
			}
		}
		if len(echoResp) == 0 {
			echoResp = append(echoResp, echo.Map{})
		}
		var breakTwoLoop bool
		for _, route := range routes {
			for _, payload := range echoResp {
				rt := RouteTask{
					Ou:          tte.Ou,
					RetrievalDB: tte.Wft.RetrievalDB,
					RouteInfo:   route,
					Payload:     payload,
				}
				fetchedResult := &hera_search.SearchResult{}
				retrievalWebTaskCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.ApiCallRequestTask, rt).Get(retrievalWebTaskCtx, &fetchedResult)
				if err != nil {
					logger.Error("failed to run api call request task retrieval", "Error", err)
					return nil, err
				}
				if fetchedResult != nil && len(fetchedResult.Value) > 0 {
					tte.Sg.SearchResults = append(tte.Sg.SearchResults, *fetchedResult)
				}
				if fetchedResult != nil && fetchedResult.WebResponse.WebFilters != nil &&
					fetchedResult.WebResponse.WebFilters.LbStrategy != nil &&
					*fetchedResult.WebResponse.WebFilters.LbStrategy != "poll-table" {
					breakTwoLoop = true
					break
				}
			}
			if breakTwoLoop {
				break
			}
		}
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update orch for retrieval services", "Error", err)
		return nil, err
	}
	return tte.Sg, nil
}

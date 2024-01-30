package ai_platform_service_orchestrations

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	lbStrategyPollTable  = "poll-table"
	lbStrategyRoundRobin = "round-robin"
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
		err = workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, tte.Ou, tte.Tc.Retrieval, tte.Sg.Window).Get(retrievalCtx, &tte.Sg)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
		tte.Sg.SourceTaskID = tte.Wft.AnalysisTaskID
	case webPlatform:
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, tte.Ou.OrgID, tte.Tc.Retrieval).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run get retrieval routes", "Error", err)
			return nil, err
		}
		for _, route := range routes {
			rt := RouteTask{
				Ou:        tte.Ou,
				Retrieval: tte.Tc.Retrieval,
				RouteInfo: route,
			}
			fetchedResult := &hera_search.SearchResult{}
			retrievalWebTaskCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.ApiCallRequestTask, rt).Get(retrievalWebTaskCtx, &fetchedResult)
			if err != nil {
				logger.Error("failed to run api call request task retrieval", "Error", err)
				return nil, err
			}
			if fetchedResult != nil && len(fetchedResult.WebResponse.Body) > 0 {
				tte.Sg.ApiResponseResults = append(tte.Sg.ApiResponseResults, *fetchedResult)
			}
			if fetchedResult != nil && fetchedResult.WebResponse.WebFilters != nil &&
				fetchedResult.WebResponse.WebFilters.LbStrategy != nil &&
				*fetchedResult.WebResponse.WebFilters.LbStrategy != lbStrategyPollTable {
				break
			}
		}
	case apiApproval:
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, tte.Ou.OrgID, tte.Wft).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
		for _, route := range routes {
			rt := RouteTask{
				Ou:        tte.Ou,
				Retrieval: tte.Tc.Retrieval,
				RouteInfo: route,
			}
			fetchedResult := &hera_search.SearchResult{}
			retrievalWebTaskCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.ApiCallRequestTask, rt).Get(retrievalWebTaskCtx, &fetchedResult)
			if err != nil {
				logger.Error("failed to run api call request task retrieval", "Error", err)
				return nil, err
			}
			if fetchedResult != nil && len(fetchedResult.WebResponse.Body) > 0 {
				tte.Sg.ApiResponseResults = append(tte.Sg.ApiResponseResults, *fetchedResult)
				trrr := &artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{
					ApprovalID:  tte.Tc.AIWorkflowTriggerResultApiResponse.ApprovalID,
					TriggerID:   tte.Tc.AIWorkflowTriggerResultApiResponse.TriggerID,
					RetrievalID: aws.ToInt(tte.Wft.RetrievalID),
					ResponseID:  tte.Tc.AIWorkflowTriggerResultApiResponse.ResponseID,
					ReqPayloads: tte.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads,
				}
				saveApiRespCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(retrievalWebTaskCtx, z.CreateOrUpdateTriggerActionApprovalWithApiReq, trrr).Get(saveApiRespCtx, &trrr)
				if err != nil {
					logger.Error("failed to save trigger response retrieval", "Error", err)
					return nil, err
				}
			}
			if fetchedResult != nil && fetchedResult.WebResponse.WebFilters != nil &&
				fetchedResult.WebResponse.WebFilters.LbStrategy != nil &&
				*fetchedResult.WebResponse.WebFilters.LbStrategy != lbStrategyPollTable {
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

func GetRetryPolicy(ret artemis_orchestrations.RetrievalItem, maxRunTime time.Duration) *temporal.RetryPolicy {
	if ret.WebFilters == nil {
		return nil
	}
	retry := &temporal.RetryPolicy{
		MaximumInterval: maxRunTime,
	}
	if ret.WebFilters.BackoffCoefficient != nil && *ret.WebFilters.BackoffCoefficient >= 1 {
		retry.BackoffCoefficient = *ret.WebFilters.BackoffCoefficient
	}
	if ret.WebFilters.MaxRetries != nil && *ret.WebFilters.MaxRetries >= 0 {
		retry.MaximumAttempts = int32(*ret.WebFilters.MaxRetries)
	}
	return retry
}

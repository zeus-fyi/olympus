package ai_platform_service_orchestrations

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	if tte.RetryPolicy != nil {
		ao.RetryPolicy = tte.RetryPolicy
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(tte.Ou.OrgID, tte.WfID, "ZeusAiPlatformServiceWorkflows", "RetrievalsWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return nil, err
	}
	platform := tte.Tc.Retrieval.RetrievalPlatform
	if tte.Tc.TriggerActionsApproval.TriggerAction == apiApproval {
		platform = apiApproval
	}

	if tte.Tc.EvalID <= 0 || tte.Tc.TriggerActionsApproval.ApprovalID <= 0 {
		platform = webPlatform
	}
	switch platform {
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
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, tte.Ou, tte.Tc.Retrieval).Get(retrievalWebCtx, &routes)
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
			apiCallCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(apiCallCtx, z.ApiCallRequestTask, rt).Get(apiCallCtx, &fetchedResult)
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
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, tte.Ou, tte.Tc.Retrieval).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}

		count := len(tte.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads)
		if count <= 0 {
			count = 1
		}
		for i := 0; i < count; i++ {
			var payload echo.Map
			if i < len(tte.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads) {
				payload = tte.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads[i]
			}
			for _, route := range routes {
				if strings.HasPrefix(route.RoutePath, "https://api.twitter.com/2") {
					if payload != nil && payload["in_reply_to_tweet_id"] != nil {
						newPayload := echo.Map{
							"reply": echo.Map{
								"in_reply_to_tweet_id": payload["in_reply_to_tweet_id"],
							},
							"text": payload["text"],
						}
						payload = newPayload
					}
				}
				rt := RouteTask{
					Ou:        tte.Ou,
					Retrieval: tte.Tc.Retrieval,
					RouteInfo: route,
					Payload:   payload,
				}
				//fmt.Println("rt", rt)
				fetchedResult := &hera_search.SearchResult{}
				apiCallCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(apiCallCtx, z.ApiCallRequestTask, rt).Get(apiCallCtx, &fetchedResult)
				if err != nil {
					logger.Error("failed to run api call request task retrieval", "Error", err)
					return nil, err
				}
				if fetchedResult != nil {
					tte.Sg.ApiResponseResults = append(tte.Sg.ApiResponseResults, *fetchedResult)
					var arrs []echo.Map
					for _, apv := range tte.Sg.ApiResponseResults {
						arrs = append(arrs, apv.WebResponse.Body)
					}
					trrr := &artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{
						ApprovalID:   tte.Tc.AIWorkflowTriggerResultApiResponse.ApprovalID,
						TriggerID:    tte.Tc.AIWorkflowTriggerResultApiResponse.TriggerID,
						RetrievalID:  aws.ToInt(tte.Tc.Retrieval.RetrievalID),
						ResponseID:   tte.Tc.AIWorkflowTriggerResultApiResponse.ResponseID,
						ReqPayloads:  tte.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads,
						RespPayloads: arrs,
					}
					bresp, berr := json.MarshalIndent(arrs, "", "  ")
					if berr != nil {
						log.Err(berr).Msg("failed to marshal resp payload")
						return nil, berr
					}

					approval := artemis_orchestrations.TriggerActionsApproval{
						TriggerAction:    apiApproval,
						ApprovalID:       tte.Tc.AIWorkflowTriggerResultApiResponse.ApprovalID,
						EvalID:           tte.Tc.EvalID,
						TriggerID:        tte.Tc.AIWorkflowTriggerResultApiResponse.TriggerID,
						WorkflowResultID: tte.Tc.TriggerActionsApproval.WorkflowResultID,
						ApprovalState:    finishedStatus,
						RequestSummary:   "Done with api call request: \n" + string(bresp),
					}
					saveApiRespCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(saveApiRespCtx, z.CreateOrUpdateTriggerActionApprovalWithApiReq, tte.Ou, approval, trrr).Get(saveApiRespCtx, &trrr)
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

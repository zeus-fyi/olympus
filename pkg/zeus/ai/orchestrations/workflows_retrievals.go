package ai_platform_service_orchestrations

import (
	"encoding/json"
	"fmt"
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

func (z *ZeusAiPlatformServiceWorkflows) RetrievalsWorkflow(ctx workflow.Context, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	if cp == nil {
		err := fmt.Errorf("wsr is nil")
		log.Err(err).Msg("RetrievalsWorkflow: failed to get workflow stage reference")
		return nil, err
	}
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    10,
		},
	}

	if cp.Tc.RetryPolicy != nil {
		ao.RetryPolicy = cp.Tc.RetryPolicy
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(cp.Ou.OrgID, cp.Wsr.ChildWfID, "ZeusAiPlatformServiceWorkflows", "RetrievalsWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return nil, err
	}
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
	switch platform {
	case twitterPlatform, redditPlatform, discordPlatform, telegramPlatform:
		retrievalCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, cp).Get(retrievalCtx, &cp.Wsr.InputID)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
	case webPlatform:
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, cp.Ou, cp.Tc.Retrieval).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run get retrieval routes", "Error", err)
			return nil, err
		}
		for _, route := range routes {
			rt := RouteTask{
				Ou:        cp.Ou,
				Retrieval: cp.Tc.Retrieval,
				RouteInfo: route,
			}
			apiCallCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(apiCallCtx, z.ApiCallRequestTask, rt).Get(apiCallCtx, &cp.Wsr.InputID)
			if err != nil {
				logger.Error("failed to run api call request task retrieval", "Error", err)
				return nil, err
			}
			if cp.Tc.Retrieval.WebFilters != nil && aws.ToString(cp.Tc.Retrieval.WebFilters.LbStrategy) != lbStrategyPollTable && cp.Wsr.InputID > 0 {
				break
			}
		}
	case apiApproval:
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, cp.Ou, cp.Tc.Retrieval).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}

		count := len(cp.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads)
		if count <= 0 {
			count = 1
		}
		for i := 0; i < count; i++ {
			var payload echo.Map
			if i < len(cp.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads) {
				payload = cp.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads[i]
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
					Ou:        cp.Ou,
					Retrieval: cp.Tc.Retrieval,
					RouteInfo: route,
					Payload:   payload,
				}
				fetchedResult := &hera_search.SearchResult{}
				apiCallCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(apiCallCtx, z.ApiCallRequestTask, rt).Get(apiCallCtx, &fetchedResult)
				if err != nil {
					logger.Error("failed to run api call request task retrieval", "Error", err)
					return nil, err
				}
				if fetchedResult != nil {
					if fetchedResult.WebResponse.Body == nil {
						fetchedResult.WebResponse.Body = echo.Map{}
					}
					trrr := &artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{
						ApprovalID:   cp.Tc.AIWorkflowTriggerResultApiResponse.ApprovalID,
						TriggerID:    cp.Tc.AIWorkflowTriggerResultApiResponse.TriggerID,
						RetrievalID:  aws.ToInt(cp.Tc.Retrieval.RetrievalID),
						ResponseID:   cp.Tc.AIWorkflowTriggerResultApiResponse.ResponseID,
						ReqPayloads:  []echo.Map{payload},
						RespPayloads: []echo.Map{fetchedResult.WebResponse.Body},
					}
					bresp, berr := json.MarshalIndent(fetchedResult.WebResponse.Body, "", "  ")
					if berr != nil {
						log.Err(berr).Msg("failed to marshal resp payload")
						return nil, berr
					}
					approval := artemis_orchestrations.TriggerActionsApproval{
						TriggerAction:    apiApproval,
						ApprovalID:       cp.Tc.AIWorkflowTriggerResultApiResponse.ApprovalID,
						EvalID:           cp.Tc.EvalID,
						TriggerID:        cp.Tc.AIWorkflowTriggerResultApiResponse.TriggerID,
						WorkflowResultID: cp.Tc.TriggerActionsApproval.WorkflowResultID,
						ApprovalState:    finishedStatus,
						RequestSummary:   "Done with api call request: \n" + string(bresp),
					}
					saveApiRespCtx := workflow.WithActivityOptions(ctx, ao)
					err = workflow.ExecuteActivity(saveApiRespCtx, z.CreateOrUpdateTriggerActionApprovalWithApiReq, cp.Ou, approval, trrr).Get(saveApiRespCtx, &trrr)
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
	return cp, nil
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

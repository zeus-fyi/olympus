package ai_platform_service_orchestrations

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
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
	ao := getRetActRetryPolicy(cp)
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(cp.Ou.OrgID, cp.Wsr.ChildWfID, "ZeusAiPlatformServiceWorkflows", "RetrievalsWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai orch services", "Error", err)
		return nil, err
	}
	switch getPlatform(cp) {
	case twitterPlatform, redditPlatform, discordPlatform, telegramPlatform:
		retrievalCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, cp).Get(retrievalCtx, &cp)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
	case webPlatform:
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, getOrgRetIfFlows(cp), cp.Tc.Retrieval).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run get retrieval routes", "Error", err)
			return nil, err
		}
		apiCallCtx := workflow.WithActivityOptions(ctx, getFanOutRetPolicy(cp, ao))
		err = workflow.ExecuteActivity(apiCallCtx, z.FanOutApiCallRequestTask, getRoutesByRule(cp, routes), cp).Get(apiCallCtx, &cp)
		if err != nil {
			logger.Error("failed to run api call request task retrieval", "Error", err)
			return nil, err
		}
	case apiApproval:
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		tmpOu := cp.Ou
		if cp.WfExecParams.WorkflowOverrides.IsUsingFlows {
			tmpOu.OrgID = FlowsOrgID
		}
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, tmpOu, cp.Tc.Retrieval).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
		count := len(cp.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads)
		if count <= 0 {
			count = 1
		} else if count > 0 && cp.Tc.Retrieval.WebFilters != nil && cp.Tc.Retrieval.WebFilters.PayloadKeys != nil {
			nem := make(map[string]bool)
			for _, key := range cp.Tc.Retrieval.WebFilters.PayloadKeys {
				nem[key] = true
			}
			for ind, pl := range cp.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads {
				for k, _ := range pl {
					if _, ok := nem[k]; !ok {
						delete(cp.Tc.AIWorkflowTriggerResultApiResponse.ReqPayloads[ind], k)
					}
				}
			}
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
				// test others first
				route.RoutePath, err = ReplaceParams(route.RoutePath, payload)
				if err != nil {
					logger.Error("failed to replace route path params", "Error", err)
					return nil, err
				}
				if cp.Tc.Retrieval.WebFilters != nil && "iterate-qp-only" == aws.ToString(cp.Tc.Retrieval.WebFilters.PayloadPreProcessing) {
					payload = nil
				}
				route.Payload = payload
				rt := RouteTask{
					Ou:        cp.Ou,
					Retrieval: cp.Tc.Retrieval,
					RouteInfo: route,
				}
				apiCallCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(apiCallCtx, z.ApiCallRequestTask, rt, cp).Get(apiCallCtx, &cp)
				if err != nil {
					logger.Error("failed to run api call request task retrieval", "Error", err)
					return nil, err
				}
				var fetchedResult hera_search.SearchResult
				if cp.Tc.ApiResponseResults != nil && len(cp.Tc.ApiResponseResults) > 0 {
					fetchedResult = cp.Tc.ApiResponseResults[0]
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
					if fetchedResult.WebResponse.WebFilters != nil &&
						fetchedResult.WebResponse.WebFilters.LbStrategy != nil &&
						*fetchedResult.WebResponse.WebFilters.LbStrategy != lbStrategyPollTable {
						break
					}
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

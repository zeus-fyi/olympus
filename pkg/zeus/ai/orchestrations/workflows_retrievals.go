package ai_platform_service_orchestrations

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
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

const (
	iterateQpOnly = "iterate-qp-only"
)

func (z *ZeusAiPlatformServiceWorkflows) RetrievalsWorkflow(ctx workflow.Context, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	if cp == nil {
		err := fmt.Errorf("wsr is nil")
		log.Err(err).Msg("RetrievalsWorkflow: failed to get workflow stage reference")
		return nil, err
	}
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 24, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			BackoffCoefficient: 2.5,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    100,
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
		err = workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, cp).Get(retrievalCtx, &cp)
		if err != nil {
			logger.Error("failed to run retrieval", "Error", err)
			return nil, err
		}
	case webPlatform:
		var routes []iris_models.RouteInfo
		retrievalWebCtx := workflow.WithActivityOptions(ctx, ao)
		tmpOu := cp.Ou
		if cp.WfExecParams.WorkflowOverrides.IsUsingFlows {
			tmpOu.OrgID = FlowsOrgID
		}
		err = workflow.ExecuteActivity(retrievalWebCtx, z.AiWebRetrievalGetRoutesTask, tmpOu, cp.Tc.Retrieval).Get(retrievalWebCtx, &routes)
		if err != nil {
			logger.Error("failed to run get retrieval routes", "Error", err)
			return nil, err
		}
		wsrCraeteCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(wsrCraeteCtx, z.CreateWsr, cp).Get(wsrCraeteCtx, &cp)
		if err != nil {
			logger.Error("failed to run get retrieval routes", "Error", err)
			return nil, err
		}
		if cp.Tc.Retrieval.RetrievalItemInstruction.WebFilters != nil &&
			cp.Tc.Retrieval.RetrievalItemInstruction.WebFilters.LbStrategy != nil &&
			*cp.Tc.Retrieval.RetrievalItemInstruction.WebFilters.LbStrategy != lbStrategyPollTable && len(routes) > 1 {
			routes = routes[0:1]
		}
		cao := ao
		cao.HeartbeatTimeout = time.Minute * 5
		cao.RetryPolicy = GetRetryPolicy(cp.Tc.Retrieval, time.Hour*24)
		cao.RetryPolicy.MaximumAttempts = 1000000
		apiCallCtx := workflow.WithActivityOptions(ctx, cao)
		err = workflow.ExecuteActivity(apiCallCtx, z.FanOutApiCallRequestTask, routes, cp).Get(apiCallCtx, &cp)
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

// ExtractParams takes a string of comma-separated regex patterns and a target string.
// It applies each regex pattern to the target string and accumulates all matched groups from each pattern into a single slice.
func ExtractParams(regexStrs []string, strContent []byte) ([]string, error) {
	// Split the regexStr into individual patterns
	var combinedParams []string
	for _, pattern := range regexStrs {
		// Compile and execute each pattern
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err // Return an error if the regular expression compilation fails
		}
		// Find all matches and extract the parameter names
		matches := re.FindAll(strContent, -1)
		for _, match := range matches {
			combinedParams = append(combinedParams, string(match))
		}
	}

	return combinedParams, nil
}

// ReplaceParams replaces placeholders in the route with URL-encoded values from the provided map.
func ReplaceParams(route string, params echo.Map) (string, error) {
	// Compile a regular expression to find {param} patterns
	re, err := regexp.Compile(`\{([^\{\}]+)\}`)
	if err != nil {
		log.Err(err).Msg("failed to compile regular expression")
		return "", err // Return an error if the regular expression compilation fails
	}

	// Replace each placeholder with the corresponding URL-encoded value from the map
	replacedRoute := re.ReplaceAllStringFunc(route, func(match string) string {
		// Extract the parameter name from the match, excluding the surrounding braces
		paramName := match[1 : len(match)-1]
		// Look up the paramName in the params map
		if value, ok := params[paramName]; ok {
			// Delete the matched entry from the map
			delete(params, paramName)
			// If the value exists, convert it to a string and URL-encode it
			return url.QueryEscape(fmt.Sprint(value))
		}
		// If no matching paramName is found in the map, return the match unchanged
		return match
	})

	return replacedRoute, nil
}

// ReplaceAndPassParams replaces placeholders in the route with URL-encoded values from the provided map.
func ReplaceAndPassParams(route string, params echo.Map) (string, []string, error) {
	// Compile a regular expression to find {param} patterns
	re, err := regexp.Compile(`\{([^\{\}]+)\}`)
	if err != nil {
		log.Err(err).Msg("failed to compile regular expression")
		return "", nil, err // Return an error if the regular expression compilation fails
	}
	var qps []string
	// Replace each placeholder with the corresponding URL-encoded value from the map
	replacedRoute := re.ReplaceAllStringFunc(route, func(match string) string {
		// Extract the parameter name from the match, excluding the surrounding braces
		paramName := match[1 : len(match)-1]
		// Look up the paramName in the params map
		if value, ok := params[paramName]; ok {
			// Delete the matched entry from the map
			if rs, rok := value.(string); rok {
				qps = append(qps, rs)
			}
			delete(params, paramName)
			// If the value exists, convert it to a string and URL-encode it
			return url.QueryEscape(fmt.Sprint(value))
		}
		// If no matching paramName is found in the map, return the match unchanged
		return match
	})

	return replacedRoute, qps, nil
}

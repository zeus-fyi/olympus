package iris_api_requests

import (
	"context"
	"errors"
	"path"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

func (i *IrisApiRequestsActivities) BroadcastETLRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	if pr == nil {
		return nil, errors.New("pr is nil")
	}
	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}
	procedureStep, ok := pr.Procedure.OrderedSteps.PopFront().(iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep)
	if !ok {
		return nil, errors.New("procedureStep not IrisRoutingProcedureStep")
	}
	payload, ok := procedureStep.BroadcastInstructions.Payload.(echo.Map)
	if !ok {
		return nil, errors.New("payload not echo.Map")
	}
	pr.Payload = payload
	pr.MaxTries = procedureStep.BroadcastInstructions.MaxTries
	// Creating a child context with a timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, procedureStep.BroadcastInstructions.MaxDuration)
	defer cancel()

	// Wait group to wait for all goroutines to complete
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var routes []iris_models.RouteInfo

	if (pr.Routes == nil || len(pr.Routes) <= 0) && (len(pr.Url) <= 0) {
		return nil, errors.New("no routes or url")
	}
	if pr.Routes != nil && len(pr.Routes) > 0 {
		copy(routes, pr.Routes)
	}
	if len(pr.Url) > 0 && (pr.Routes == nil || len(pr.Routes) == 0) {
		routes = append(routes, iris_models.RouteInfo{RoutePath: pr.Url, Referers: pr.Referrers})
	}
	if len(procedureStep.BroadcastInstructions.RoutingPath) > 0 {
		pr.ExtRoutePath = procedureStep.BroadcastInstructions.RoutingPath
	}
	if len(pr.ExtRoutePath) > 0 || len(pr.Referrers) > 0 {
		for ind, _ := range routes {
			routes[ind].RoutePath = path.Join(routes[ind].RoutePath, pr.ExtRoutePath)
			routes[ind].Referers = pr.Referrers
		}
	}
	// Iterating through routes and launching goroutines
	for _, route := range routes {
		wg.Add(1)
		go func(ctx context.Context, r string, cancel func()) {
			defer wg.Done()
			// Make a copy of the ApiProxyRequest to avoid race conditions
			req := *pr
			req.Url = r

			// Call ExtLoadBalancerRequest with the modified request
			resp, err := i.ExtLoadBalancerRequest(ctx, &req)
			if err == nil && resp.StatusCode < 400 {
				for _, transform := range procedureStep.TransformSlice {
					// Assuming that resp.Response contains the data from which to extract the key value
					transform.Source = r
					transform.ExtractKeyValue(resp.Response)
					mutex.Lock() // Lock access to shared procedureStep
					pr.PayloadSizeMeter.Add(resp.PayloadSizeMeter.Size)
					if len(transform.ExtractionKey) > 0 {
						agg, aok := procedureStep.AggregateMap[transform.ExtractionKey]
						if aok {
							aerr := agg.AggregateOn(transform.Value, transform)
							if aerr != nil {
								log.Err(aerr).Msg("failed to aggregate")
								return
							}
							procedureStep.AggregateMap[transform.ExtractionKey] = agg
						}
					}
					if procedureStep.BroadcastInstructions.FanInRules != nil {
						switch procedureStep.BroadcastInstructions.FanInRules.Rule {
						case iris_programmable_proxy_v1_beta.FanInRuleFirstValidResponse:
							pr = resp
							cancel()
						}
					}
					mutex.Unlock()
				}
			} else {
				mutex.Lock()
				pr.PayloadSizeMeter.Add(resp.PayloadSizeMeter.Size)
				mutex.Unlock()
				log.Err(err).Msg("Failed to broadcast request")
			}
		}(timeoutCtx, route.RoutePath, cancel)
	}
	// Wait for all goroutines to complete
	wg.Wait()
	if pr.Procedure.OrderedSteps.Len() > 0 {
		return i.BroadcastETLRequest(ctx, pr)
	}
	return pr, nil
}

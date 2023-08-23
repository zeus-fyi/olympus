package iris_api_requests

import (
	"context"
	"errors"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func (i *IrisApiRequestsActivities) BroadcastETLRequest(ctx context.Context, pr *ApiProxyRequest, routes []iris_models.RouteInfo) (*ApiProxyRequest, error) {
	if len(pr.Procedure.OrderedSteps) == 0 {
		return nil, errors.New("no steps in procedure")
	}

	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}

	procedureStep := pr.Procedure.OrderedSteps[0]
	payload, ok := procedureStep.BroadcastInstructions.Payload.(echo.Map)
	if !ok {
		return nil, errors.New("payload not echo.Map")
	}
	pr.Payload = payload
	// Creating a child context with a timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, procedureStep.BroadcastInstructions.MaxDuration)
	defer cancel()

	// Wait group to wait for all goroutines to complete
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Iterating through routes and launching goroutines
	for _, route := range routes {
		wg.Add(1)
		go func(ctx context.Context, r string) {
			defer wg.Done()

			// Make a copy of the ApiProxyRequest to avoid race conditions
			req := *pr
			req.Url = r

			// Call ExtLoadBalancerRequest with the modified request
			resp, err := i.ExtLoadBalancerRequest(ctx, &req)
			if err == nil {
				for _, transform := range procedureStep.TransformSlice {
					// Assuming that resp.Response contains the data from which to extract the key value
					transform.Source = r
					transform.ExtractKeyValue(resp.Response)
					mutex.Lock() // Lock access to shared procedureStep
					pr.PayloadSizeMeter.Add(resp.PayloadSizeMeter.Size)
					agg, aok := procedureStep.AggregateMap[transform.ExtractionKey]
					if aok {
						aerr := agg.AggregateOn(transform.Value, transform)
						if aerr != nil {
							log.Err(aerr).Msg("Failed to aggregate")
						}
						procedureStep.AggregateMap[transform.ExtractionKey] = agg
					}
					mutex.Unlock() // Unlock access to shared procedureStep
				}
			} else {
				mutex.Lock()
				pr.PayloadSizeMeter.Add(resp.PayloadSizeMeter.Size)
				mutex.Unlock()
				log.Err(err).Msg("Failed to broadcast request")
			}
		}(timeoutCtx, route.RoutePath)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return pr, nil
}

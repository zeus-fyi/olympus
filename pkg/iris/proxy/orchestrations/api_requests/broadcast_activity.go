package iris_api_requests

import (
	"context"
	"errors"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (i *IrisApiRequestsActivities) BroadcastETLRequest(ctx context.Context, pr *ApiProxyRequest, routes []iris_models.RouteInfo) (*ApiProxyRequest, error) {
	if len(pr.Procedure.OrderedSteps) == 0 {
		return nil, errors.New("no steps in procedure")
	}

	procedureStep := pr.Procedure.OrderedSteps[0]

	payload, ok := procedureStep.BroadcastInstructions.Payload.(echo.Map)
	if !ok {
		return nil, errors.New("payload not echo.Map")
	}
	pr.Payload = payload
	// Creating a child context with a timeout
	//timeoutCtx, cancel := context.WithTimeout(ctx, procedureStep.BroadcastInstructions.MaxDuration)
	//defer cancel()

	// Channel to collect the results
	results := make(chan *ApiProxyRequest, len(routes))

	// Wait group to wait for all goroutines to complete
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Iterating through routes and launching goroutines
	for _, route := range routes {
		wg.Add(1)
		go func(r string) {
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
				results <- resp
			}
		}(route.RoutePath)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Close the channel to stop the receiver
	close(results)

	// Process the results as needed
	var finalResponse *ApiProxyRequest
	// You can choose how to aggregate or select the final response
	for result := range results {
		finalResponse = result
		// Additional logic to combine or select responses
	}

	return finalResponse, nil
}

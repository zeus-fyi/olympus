package iris_api_requests

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

type IrisApiRequestsActivities struct {
}

func NewArtemisApiRequestsActivities() IrisApiRequestsActivities {
	return IrisApiRequestsActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (i *IrisApiRequestsActivities) GetActivities() ActivitiesSlice {
	return []interface{}{i.RelayRequest, i.InternalSvcRelayRequest, i.ExtLoadBalancerRequest, i.UpdateOrgRoutingTable,
		i.SelectSingleOrgGroupsRoutingTables, i.SelectOrgGroupRoutingTable, i.SelectAllRoutingTables,
		i.DeleteOrgRoutingTable,
	}
}

func (i *IrisApiRequestsActivities) RelayRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)
	if pr.IsInternal {
		r.SetAuthToken(artemis_orchestration_auth.Bearer)
	}
	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}
	resp, err := r.R().SetBody(&pr.Payload).SetResult(&pr.Response).Post(pr.Url)
	if err != nil {
		log.Err(err).Msg("Failed to relay api request")
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Msg("Failed to relay api request")
		return nil, fmt.Errorf("failed to relay api request: status code %d", resp.StatusCode())
	}
	return pr, err
}

func (i *IrisApiRequestsActivities) InternalSvcRelayRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)
	if pr.IsInternal {
		r.SetAuthToken(artemis_orchestration_auth.Bearer)
	}
	resp, err := r.R().SetBody(&pr.Payload).SetResult(&pr.Response).Post(pr.Url)
	if err != nil {
		log.Err(err).Interface("statusCode", resp.StatusCode()).Msg("Failed to relay api request")
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("statusCode", resp.StatusCode()).Msg("Failed to relay api request")
		return nil, fmt.Errorf("failed to relay api request: status code %d", resp.StatusCode())
	}
	return pr, err
}

func (i *IrisApiRequestsActivities) BroadcastETLRequest(ctx context.Context, pr *ApiProxyRequest, routes []string, timeout time.Duration) (*ApiProxyRequest, error) {
	// Creating a child context with a timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Channel to collect the results
	results := make(chan *ApiProxyRequest, len(routes))

	// Wait group to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Iterating through routes and launching goroutines
	for _, route := range routes {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()

			// Make a copy of the ApiProxyRequest to avoid race conditions
			req := *pr
			req.Url = r

			// Call ExtLoadBalancerRequest with the modified request
			resp, err := i.ExtLoadBalancerRequest(timeoutCtx, &req)
			if err == nil {
				results <- resp
			}
		}(route)
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

func (i *IrisApiRequestsActivities) ExtLoadBalancerRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)

	parsedURL, err := url.Parse(pr.Url)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return pr, err
	}

	if parsedURL.Scheme != "https" {
		return pr, fmt.Errorf("error: URL must be an HTTPS URL")
	}

	if len(pr.QueryParams) > 0 {
		r.QueryParam = pr.QueryParams
	}

	for k, v := range pr.RequestHeaders {
		r.SetHeader(k, strings.Join(v, ", ")) // Joining all values with a comma
	}

	if len(pr.Referrers) > 0 {
		r.SetHeader("Referer", strings.Join(pr.Referrers, ", ")) // Joining all values with a comma
	}

	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}
	var resp *resty.Response
	switch pr.PayloadTypeREST {
	case "GET":
		resp, err = sendRequest(r.R(), pr, "GET")
	case "PUT":
		resp, err = sendRequest(r.R(), pr, "PUT")
	case "DELETE":
		resp, err = sendRequest(r.R(), pr, "DELETE")
	default:
		resp, err = sendRequest(r.R(), pr, "POST")
	}
	if err != nil {
		log.Err(err).Msg("Failed to relay api request")
		return pr, err
	}
	if pr.StatusCode >= 400 {
		log.Err(err).Msg("Failed to relay api request")
		return pr, fmt.Errorf("failed to relay api request: status code %d", resp.StatusCode())
	}
	return pr, err
}

func sendRequest(request *resty.Request, pr *ApiProxyRequest, method string) (*resty.Response, error) {
	var resp *resty.Response
	var err error

	if pr.Payload != nil {
		switch method {
		case "GET":
			resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Get(pr.Url)
		case "PUT":
			resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Put(pr.Url)
		case "DELETE":
			resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Delete(pr.Url)
		default:
			resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Post(pr.Url)
		}
	} else {
		switch method {
		case "GET":
			resp, err = request.SetResult(&pr.Response).Get(pr.Url)
		case "PUT":
			resp, err = request.SetResult(&pr.Response).Put(pr.Url)
		case "DELETE":
			resp, err = request.SetResult(&pr.Response).Delete(pr.Url)
		default:
			resp, err = request.SetResult(&pr.Response).Post(pr.Url)
		}
	}
	if resp != nil {
		pr.PayloadSizeMeter.Add(resp.Size())
		pr.StatusCode = resp.StatusCode()
		if resp.StatusCode() >= 400 || pr.Response == nil {
			pr.RawResponse = resp.Body()
		}
		pr.ResponseHeaders = filterHeaders(resp.RawResponse.Header)

		pr.ReceivedAt = resp.ReceivedAt()
		pr.Latency = resp.Time()
	}
	return resp, err
}

func filterHeaders(headers http.Header) http.Header {
	filteredHeaders := make(http.Header)
	for key, values := range headers {
		if key[:2] == "X-" {
			filteredHeaders[key] = values
		}
	}
	return filteredHeaders
}

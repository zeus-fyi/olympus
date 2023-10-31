package iris_api_requests

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

type IrisApiRequestsActivities struct {
}

func NewIrisApiRequestsActivities() IrisApiRequestsActivities {
	return IrisApiRequestsActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (i *IrisApiRequestsActivities) GetActivities() ActivitiesSlice {
	return []interface{}{i.RelayRequest, i.InternalSvcRelayRequest, i.ExtLoadBalancerRequest, i.UpdateOrgRoutingTable,
		i.SelectSingleOrgGroupsRoutingTables, i.SelectOrgGroupRoutingTable, i.SelectAllRoutingTables,
		i.DeleteOrgRoutingTable, i.ExtToAnvilInternalSimForkRequest,
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

func (i *IrisApiRequestsActivities) ExtToAnvilInternalSimForkRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	pr.IsInternal = true
	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}
	return i.ExtLoadBalancerRequest(ctx, pr)
}

func (i *IrisApiRequestsActivities) ExtLoadBalancerRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)
	if pr.MaxTries > 0 {
		r.SetRetryCount(pr.MaxTries)
	}
	parsedURL, err := url.Parse(pr.Url)
	if err != nil {
		log.Err(err).Msg("ExtLoadBalancerRequest: failed to parse url")
		return pr, err
	}

	if pr.OrgID == 7138983863666903883 {
		// for internal
	} else if pr.IsInternal {
		log.Info().Interface("pr.URL", pr.Url).Msg("ExtLoadBalancerRequest: anvil request")
	} else {
		if parsedURL.Scheme != "https" {
			err = fmt.Errorf("error: URL must be an HTTPS URL")
			log.Err(err).Msg("ExtLoadBalancerRequest: http request unauthorized")
			return pr, err
		}
	}

	if len(pr.QueryParams) > 0 {
		r.QueryParam = pr.QueryParams
	}

	for k, v := range pr.RequestHeaders {
		switch k {
		case "Authorization":
			continue
		}
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
	case "POST":
		resp, err = sendRequest(r.R(), pr, "POST")
	default:
		resp, err = sendRequest(r.R(), pr, "POST")
	}
	if err != nil {
		log.Err(err).Msg("ExtLoadBalancerRequest: Failed to relay api request")
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

	ext := ""
	if pr.ExtRoutePath != "" {
		ext = pr.ExtRoutePath
	}
	if pr.Payload != nil {
		switch method {
		case "GET":
			resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Get(ext)
		case "PUT":
			resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Put(ext)
		case "DELETE":
			resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Delete(ext)
		case "POST":
			resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Post(ext)
		default:
			resp, err = request.SetBody(&pr.Payload).SetResult(&pr.Response).Post(ext)
		}
	} else {
		switch method {
		case "GET":
			resp, err = request.SetResult(&pr.Response).Get(ext)
		case "PUT":
			resp, err = request.SetResult(&pr.Response).Put(ext)
		case "DELETE":
			resp, err = request.SetResult(&pr.Response).Delete(ext)
		case "POST":
			resp, err = request.SetResult(&pr.Response).Post(ext)
		default:
			resp, err = request.SetResult(&pr.Response).Post(ext)
		}
	}
	if err != nil {
		log.Err(err).Msg("Failed to relay api request")
		return nil, err
	}

	if resp != nil {
		if pr.PayloadSizeMeter == nil {
			pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
		}
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

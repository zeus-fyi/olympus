package iris_api_requests

import (
	"context"
	"fmt"

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

func (i *IrisApiRequestsActivities) ExtLoadBalancerRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)
	// only get first referer, not sure why a list is needed
	for ind, ref := range pr.Referrers {
		if ind == 0 {
			r.SetHeader("Referer", ref)
		}
	}
	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}
	var resp *resty.Response
	var err error
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

	pr.PayloadSizeMeter.Add(resp.Size())
	if resp.StatusCode() >= 400 {
		log.Err(err).Msg("Failed to relay api request")
		return nil, fmt.Errorf("failed to relay api request: status code %d", resp.StatusCode())
	}
	return pr, err
}

func sendRequest(request *resty.Request, pr *ApiProxyRequest, method string) (*resty.Response, error) {
	var resp *resty.Response
	var err error

	if pr.PayloadSizeMeter.N() > int64(0) {
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
	return resp, err
}

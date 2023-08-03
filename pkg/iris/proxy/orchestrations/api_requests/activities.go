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
		i.DeleteOrgRoutingTable, i.ExtLoadBalancerGETRequest,
	}
}

func (i *IrisApiRequestsActivities) RelayRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)
	if pr.IsInternal {
		r.SetAuthToken(artemis_orchestration_auth.Bearer)
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
	switch pr.PayloadTypeREST {
	case "GET":
		return i.ExtLoadBalancerGETRequest(ctx, pr)
	}
	r := resty.New()
	r.SetBaseURL(pr.Url)
	// only get first referer, not sure why a list is needed
	for ind, ref := range pr.Referrers {
		if ind == 0 {
			r.SetHeader("Referer", ref)
		}
	}
	resp, err := r.R().SetBody(&pr.Payload).SetResult(&pr.Response).Post(pr.Url)
	if err != nil {
		log.Err(err).Msg("Failed to relay api request")
		return nil, err
	}
	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}
	pr.PayloadSizeMeter.Add(resp.Size())
	if resp.StatusCode() >= 400 {
		log.Err(err).Msg("Failed to relay api request")
		return nil, fmt.Errorf("failed to relay api request: status code %d", resp.StatusCode())
	}
	return pr, err
}

func (i *IrisApiRequestsActivities) ExtLoadBalancerGETRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)
	// only get first referer, not sure why a list is needed
	for ind, ref := range pr.Referrers {
		if ind == 0 {
			r.SetHeader("Referer", ref)
		}
	}
	resp, err := r.R().SetResult(&pr.Response).Get(pr.Url)
	if err != nil {
		log.Err(err).Msg("Failed to relay api request")
		return nil, err
	}
	if pr.PayloadSizeMeter == nil {
		pr.PayloadSizeMeter = &iris_usage_meters.PayloadSizeMeter{}
	}
	pr.PayloadSizeMeter.Add(resp.Size())
	if resp.StatusCode() >= 400 {
		log.Err(err).Msg("Failed to relay api request")
		return nil, fmt.Errorf("failed to relay api request: status code %d", resp.StatusCode())
	}
	return pr, err
}

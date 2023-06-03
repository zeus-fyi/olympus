package artemis_api_requests

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type ArtemisApiRequestsActivities struct {
}

func NewArtemisApiRequestsActivities() ArtemisApiRequestsActivities {
	return ArtemisApiRequestsActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (a *ArtemisApiRequestsActivities) GetActivities() ActivitiesSlice {
	return []interface{}{a.RelayRequest}
}

func (a *ArtemisApiRequestsActivities) RelayRequest(ctx context.Context, pr ApiProxyRequest) ([]byte, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)
	resp, err := r.R().SetBody(&pr.Payload).Post(pr.Url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), err
}

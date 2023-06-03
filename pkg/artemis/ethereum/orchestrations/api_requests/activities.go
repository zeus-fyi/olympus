package artemis_api_requests

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
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

func (a *ArtemisApiRequestsActivities) RelayRequest(ctx context.Context, pr *ApiProxyRequest) (*ApiProxyRequest, error) {
	r := resty.New()
	r.SetBaseURL(pr.Url)
	if pr.IsInternal {
		r.SetAuthToken(artemis_orchestration_auth.Bearer)
	}
	_, err := r.R().SetBody(&pr.Payload).SetResult(&pr.Response).Post(pr.Url)
	if err != nil {
		log.Err(err).Msg("Failed to relay api request")
		return nil, err
	}
	return pr, err
}

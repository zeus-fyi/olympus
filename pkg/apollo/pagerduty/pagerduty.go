package apollo_pagerduty

import (
	"context"
	"errors"
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

const (
	PagerDutyAPI = "https://events.pagerduty.com/v2/enqueue"
)

type PagerDutyClient struct {
	EventAction
	Severity
	*resty.Client
}

func NewPagerDutyClient(apiKey string) PagerDutyClient {
	return PagerDutyClient{
		Client: resty_base.GetBaseRestyClient(PagerDutyAPI, apiKey).Client,
	}
}

type V2EventResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	DedupKey string `json:"dedup_key"`
}

func (pd *PagerDutyClient) SendAlert(ctx context.Context, event pagerduty.V2Event) (V2EventResponse, error) {
	respJSON := V2EventResponse{}
	resp, err := pd.R().
		SetResult(&respJSON).
		SetBody(event).
		Post("/")

	if err != nil || resp.StatusCode() >= 400 {
		log.Err(err).Msg("PagerDutyClient: SendAlert")
		if resp.StatusCode() >= 400 {
			err = errors.New(fmt.Sprintf("bad request: code: %d", resp.StatusCode()))
			log.Err(err).Interface("event", event).Msg("PagerDutyClient: SendAlert")
		}
		if err == nil {
			err = fmt.Errorf("non-OK status code: %d", resp.StatusCode())
		}
		return respJSON, err
	}
	return respJSON, err
}

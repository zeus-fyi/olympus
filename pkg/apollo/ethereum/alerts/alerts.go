package apollo_ethereum_alerts

import (
	"context"
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/rs/zerolog/log"
	apollo_pagerduty "github.com/zeus-fyi/olympus/pkg/apollo/pagerduty"
	apollo_prometheus "github.com/zeus-fyi/olympus/pkg/apollo/prometheus"
)

type ApolloEthereumAlerts struct {
	pagerdutyRoutingKey string
	apollo_pagerduty.PagerDutyClient
	apollo_prometheus.Prometheus
}

func InitLocalApolloEthereumAlerts(ctx context.Context, pdApiKey, pagerdutyRoutingKey string) ApolloEthereumAlerts {
	return ApolloEthereumAlerts{pagerdutyRoutingKey, apollo_pagerduty.NewPagerDutyClient(pdApiKey), apollo_prometheus.NewPrometheusLocalClient(ctx)}
}

func InitProdApolloEthereumAlerts(ctx context.Context, pdApiKey, pagerdutyRoutingKey string) ApolloEthereumAlerts {
	return ApolloEthereumAlerts{pagerdutyRoutingKey, apollo_pagerduty.NewPagerDutyClient(pdApiKey), apollo_prometheus.NewPrometheusProdClient(ctx)}
}

func (a *ApolloEthereumAlerts) CreateAlertFromTemplate(ctx context.Context, summary, component, details string) (apollo_pagerduty.V2EventResponse, error) {
	event := pagerduty.V2Event{
		RoutingKey: a.pagerdutyRoutingKey,
		Action:     a.Trigger(),
		Payload: &pagerduty.V2Payload{
			Summary:   summary,
			Source:    "APOLLO_ETHEREUM_ALERTS",
			Severity:  a.Critical(),
			Component: fmt.Sprintf("This is the %s component", component),
			Group:     "This is the ethereum group",
			Class:     "Ethereum",
			Details:   details,
		},
	}
	r, err := a.SendAlert(ctx, event)
	if err != nil {
		log.Error().Err(err).Msg("apollo_ethereum_alerts.CreateAlertFromTemplate")
		return r, err
	}
	return r, err
}

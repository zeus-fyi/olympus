package apollo_prometheus

import (
	"context"

	"github.com/prometheus/client_golang/api"
	"github.com/rs/zerolog/log"
)

const (
	localPrometheusHostPort = "http://localhost:9090"
	prodPrometheusHostPort  = "http://prometheus-operated.observability.svc.cluster.local:9090"
)

type Prometheus struct {
	api.Client
}

func NewPrometheusClient(ctx context.Context, hostURL string) Prometheus {
	promClient, err := api.NewClient(api.Config{Address: hostURL})
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("apollo_prometheus.NewPrometheusClient")
		panic(err)
	}
	return Prometheus{promClient}
}

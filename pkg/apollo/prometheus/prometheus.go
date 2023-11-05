package apollo_prometheus

import (
	"context"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/rs/zerolog/log"
)

const (
	localPrometheusHostPort = "http://localhost:9090"
	prodPrometheusHostPort  = "http://prometheus-operated.observability.svc.cluster.local:9090"
)

var ProxyPromQL Prometheus

type Prometheus struct {
	v1.API
}

func NewPrometheusLocalClient(ctx context.Context) Prometheus {
	return NewPrometheusClient(ctx, localPrometheusHostPort)
}
func NewPrometheusProdClient(ctx context.Context) Prometheus {
	return NewPrometheusClient(ctx, prodPrometheusHostPort)
}

func NewPrometheusClient(ctx context.Context, hostURL string) Prometheus {
	promClient, err := api.NewClient(api.Config{Address: hostURL})
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("apollo_prometheus.NewPrometheusClient")
		panic(err)
	}
	apiClient := v1.NewAPI(promClient)
	return Prometheus{apiClient}
}

func InitNewPromQLProdClient(ctx context.Context) {
	ProxyPromQL = NewPrometheusClient(ctx, prodPrometheusHostPort)
}

func InitProdLocalPromQLProdClient(ctx context.Context) {
	ProxyPromQL = NewPrometheusClient(ctx, localPrometheusHostPort)
}

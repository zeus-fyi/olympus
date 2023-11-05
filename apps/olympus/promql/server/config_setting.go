package promql_server

import (
	"context"

	mev_promql "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/promql"
	apollo_prometheus "github.com/zeus-fyi/olympus/pkg/apollo/prometheus"
)

func SetConfigByEnv(ctx context.Context, env string) {
	switch env {
	case "production":
		apollo_prometheus.InitNewPromQLProdClient(ctx)
	case "production-local":
		apollo_prometheus.InitProdLocalPromQLProdClient(ctx)
	case "local":
		apollo_prometheus.InitProdLocalPromQLProdClient(ctx)
	}
	mev_promql.ProxyMevPromQL = mev_promql.NewMevPromQL(apollo_prometheus.ProxyPromQL)
}

package hyrda_choreography_metrics

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	apollo_metrics "github.com/zeus-fyi/olympus/pkg/apollo/metrics"
	apollo_beacon_prom_metrics "github.com/zeus-fyi/olympus/pkg/apollo/metrics/ethereum/beacon"
	apollo_metrics_workload_info "github.com/zeus-fyi/olympus/pkg/apollo/metrics/workload_info"
)

var HydraChoreographyMetrics *prometheus.Registry

func MetricsRequestHandler(c echo.Context) error {
	promHandler := promhttp.HandlerFor(HydraChoreographyMetrics, promhttp.HandlerOpts{Registry: HydraChoreographyMetrics})
	return echo.WrapHandler(promHandler)(c)
}

func InitHydraMetrics(ctx context.Context, wi apollo_metrics_workload_info.WorkloadInfo) {
	HydraChoreographyMetrics = apollo_metrics.InitMetrics()
	InitBeaconMetricsForMonitoring(ctx, wi)
}

func InitBeaconMetricsForMonitoring(ctx context.Context, wi apollo_metrics_workload_info.WorkloadInfo) {
	bc := apollo_beacon_prom_metrics.HydraConfig("ephemeral")
	bm := apollo_beacon_prom_metrics.NewBeaconMetrics(wi, bc, "")
	HydraChoreographyMetrics.MustRegister(bm.GetMetrics()...)

	go bm.PollMetrics(10 * time.Second)
}

package tyche_metrics

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	apollo_metrics "github.com/zeus-fyi/olympus/pkg/apollo/metrics"
	apollo_metrics_workload_info "github.com/zeus-fyi/olympus/pkg/apollo/metrics/workload_info"
)

var TycheMetrics *prometheus.Registry

func MetricsRequestHandler(c echo.Context) error {
	promHandler := promhttp.HandlerFor(TycheMetrics, promhttp.HandlerOpts{Registry: TycheMetrics})
	return echo.WrapHandler(promHandler)(c)
}

func InitTycheMetrics(ctx context.Context, wi apollo_metrics_workload_info.WorkloadInfo) {
	TycheMetrics = apollo_metrics.InitMetrics()
}

func InitBeaconMetricsForMonitoring(ctx context.Context, wi apollo_metrics_workload_info.WorkloadInfo) {
	//	bc := apollo_beacon_prom_metrics.HydraConfig("ephemeral")
	//	bm := apollo_beacon_prom_metrics.NewBeaconMetrics(wi, bc, "")
	//	TycheMetrics.MustRegister(bm.GetMetrics()...)
	//
	//	go bm.PollMetrics(10 * time.Second)
}

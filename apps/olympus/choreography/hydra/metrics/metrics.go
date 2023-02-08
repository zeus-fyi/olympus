package hyrda_choreography_metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	apollo_metrics "github.com/zeus-fyi/olympus/pkg/apollo/metrics"
	apollo_beacon_prom_metrics "github.com/zeus-fyi/olympus/pkg/apollo/metrics/ethereum/beacon"
	apollo_metrics_workload_info "github.com/zeus-fyi/olympus/pkg/apollo/metrics/workload_info"
	"time"
)

var HydraChoreographyMetrics *prometheus.Registry

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

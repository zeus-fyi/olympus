package apollo_beacon_prom_metrics

import (
	"github.com/stretchr/testify/suite"
	apollo_metrics_workload_info "github.com/zeus-fyi/olympus/pkg/apollo/metrics/workload_info"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	"testing"
)

type BeaconMetricsClientTestSuite struct {
	test_suites_base.TestSuite
	BeaconMetricsClient BeaconMetrics
}

func (t *BeaconMetricsClientTestSuite) SetupTest() {
	t.InitLocalConfigs()
	wi := apollo_metrics_workload_info.WorkloadInfo{}
	t.BeaconMetricsClient = NewBeaconMetrics(wi, BeaconConfig{}, t.Tc.ProductionLocalTemporalBearerToken)
}

func TestBeaconMetricsClientTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconMetricsClientTestSuite))
}

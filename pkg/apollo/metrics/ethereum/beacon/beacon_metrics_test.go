package apollo_beacon_prom_metrics

import (
	"testing"

	"github.com/stretchr/testify/suite"
	apollo_metrics_workload_info "github.com/zeus-fyi/olympus/pkg/apollo/metrics/workload_info"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	client_consts "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons/constants"
	resty_base "github.com/zeus-fyi/zeus/pkg/zeus/client/base"
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

func (t *BeaconMetricsClientTestSuite) TestBeaconExecClientSyncStatus() {
	ss := client_consts.ExecClientSyncStatus{}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	r := resty_base.GetBaseRestyClient(t.Tc.EphemeralNodeUrl, "")
	resp, err := r.R().
		SetHeaders(headers).
		SetResult(&ss).
		SetBody(beaconExecSyncPayload).Post("/")

	t.Require().NoError(err)
	t.Require().Equal(200, resp.StatusCode())
	t.Require().NotEmpty(ss)

	t.Assert().False(ss.Result)
}

func TestBeaconMetricsClientTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconMetricsClientTestSuite))
}

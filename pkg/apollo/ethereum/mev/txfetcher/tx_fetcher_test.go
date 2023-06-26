package metrics_txfetcher

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type TxFetcherMetricsTestSuite struct {
	test_suites_base.TestSuite
	m TxFetcherMetrics
}

func (t *TxFetcherMetricsTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func TestTxFetcherMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(TxFetcherMetricsTestSuite))
}

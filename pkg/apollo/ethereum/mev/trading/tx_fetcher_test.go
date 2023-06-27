package metrics_trading

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type TxFetcherMetricsTestSuite struct {
	test_suites_base.TestSuite
	m TxFetcherMetrics
}

func (t *TxFetcherMetricsTestSuite) TestNewTxFetcherMetrics() {
	// Create a new registry to register the metrics
	reg := prometheus.NewPedanticRegistry()

	// Create a new TxFetcherMetrics and register it to the local registry
	txMetrics := NewTxFetcherMetrics(reg)

	// Simulate a transaction
	txMetrics.TransactionGroup(web3_client.UniswapUniversalRouterAddressNew, "swapExactETHForTokens")

	// Use the GatherAndCount helper function to count the number of occurrences of the metric
	count, err := testutil.GatherAndCount(reg, "ethereum_mev_tx_stats")
	t.Require().NoError(err)

	// Assert that the count is 1
	t.Equal(1, count)
}

func (t *TxFetcherMetricsTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func TestTxFetcherMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(TxFetcherMetricsTestSuite))
}

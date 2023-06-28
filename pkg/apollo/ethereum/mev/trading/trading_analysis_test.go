package metrics_trading

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type TradeAnalysisMetricsTestSuite struct {
	test_suites_base.TestSuite
	m TradingMetrics
}

func (t *TradeAnalysisMetricsTestSuite) TestTradeAnalysisMetrics() {
	// Create a new registry to register the metrics
	reg := prometheus.NewPedanticRegistry()
	// Create a new TxFetcherMetrics and register it to the local registry
	txMetrics := NewTradingMetrics(reg)
	// The method and token
	method := "method1"
	token := "token1"
	revenue := 100.0
	cost := 10.0
	txMetrics.TradeAnalysisMetrics.CalculatedSandwich(method, token, revenue, cost)
	// Use the GatherAndCount helper function to count the number of occurrences of the metric
	count, err := testutil.GatherAndCount(reg, "eth_mev_sandwich_calculated_revenue")
	t.Require().NoError(err)
	// Assert that the count is 1
	t.Equal(1, count)
	count, err = testutil.GatherAndCount(reg, "eth_mev_sandwich_calculated_revenue_event")
	t.Require().NoError(err)
	t.Equal(1, count)
	count, err = testutil.GatherAndCount(reg, "eth_mev_sandwich_upfront_cost_event")
	t.Require().NoError(err)
	t.Equal(1, count)
	count, err = testutil.GatherAndCount(reg, "eth_mev_sandwich_calculated_roi_event")
	t.Require().NoError(err)
	t.Equal(1, count)
}

func (t *TradeAnalysisMetricsTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func TestTradeAnalysisMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(TradeAnalysisMetricsTestSuite))
}

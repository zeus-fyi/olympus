package metrics_trading

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type TradeAnalysisMetricsTestSuite struct {
	test_suites_base.TestSuite
	m TradingMetrics
}

func (t *TradeAnalysisMetricsTestSuite) TestTradeAnalysisMetrics() {
	// Create a new registry to register the metrics
	//reg := prometheus.NewPedanticRegistry()
}

func (t *TradeAnalysisMetricsTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func TestTradeAnalysisMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(TradeAnalysisMetricsTestSuite))
}

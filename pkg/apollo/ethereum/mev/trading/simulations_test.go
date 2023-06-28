package metrics_trading

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type TradeSimulationMetricsTestSuite struct {
	test_suites_base.TestSuite
	m TradingMetrics
}

func (t *TradeSimulationMetricsTestSuite) TestTradeSimulationMetrics() {
	// Create a new registry to register the metrics
	//reg := prometheus.NewPedanticRegistry()
}

func (t *TradeSimulationMetricsTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func TestTradeSimulationMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(TradeSimulationMetricsTestSuite))
}

package metrics_simulations

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type SimulationMetricsTestSuite struct {
	test_suites_base.TestSuite
	s SimulationMetrics
}

func (t *SimulationMetricsTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func TestSimulationMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(SimulationMetricsTestSuite))
}

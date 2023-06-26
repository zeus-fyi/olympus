package metrics_flashbots

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type FlashbotsMetricsTestSuite struct {
	test_suites_base.TestSuite
	f FlashbotsMetrics
}

func (t *FlashbotsMetricsTestSuite) SetupTest() {
	t.InitLocalConfigs()
}

func TestFlashbotsMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(FlashbotsMetricsTestSuite))
}

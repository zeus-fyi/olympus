package apollo_prometheus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type PrometheusTestSuite struct {
	pc Prometheus
	test_suites_base.TestSuite
}

func (t *PrometheusTestSuite) SetupTest() {
	t.InitLocalConfigs()
	t.pc = NewPrometheusClient(ctx, localPrometheusHostPort)
}

func (t *PrometheusTestSuite) TestPromQuery() {

	t.Require().NotEmpty(t.pc)

}

func TestPrometheusTestSuite(t *testing.T) {
	suite.Run(t, new(PrometheusTestSuite))
}

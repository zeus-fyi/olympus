package mev_promql

import (
	"context"
	"fmt"
	"testing"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/stretchr/testify/suite"
	apollo_prometheus "github.com/zeus-fyi/olympus/pkg/apollo/prometheus"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type MevPrometheusTestSuite struct {
	pc MevPromQL
	test_suites_base.TestSuite
}

func (t *MevPrometheusTestSuite) SetupTest() {
	t.InitLocalConfigs()
	t.pc = NewMevPromQL(apollo_prometheus.NewPrometheusLocalClient(ctx))
	t.pc.printOn = true
}

func (t *MevPrometheusTestSuite) TestQueryRangePromQL() {
	t.Require().NotEmpty(t.pc)

	timeNow := time.Now().UTC()

	window := v1.Range{
		Start: timeNow.Add(-time.Minute * 60),
		End:   time.Now().UTC(),
		Step:  time.Minute,
	}
	t.pc.printOn = false
	m, err := t.pc.GetTopRevenuePairs(ctx, window)
	fmt.Println(m)
	t.Require().NoError(err)
	t.Assert().NotEmpty(m)

	for _, val := range m {
		fmt.Println(val.Metric.In, val.Metric.Pair)
	}
}

func TestMevPrometheusTestSuite(t *testing.T) {
	suite.Run(t, new(MevPrometheusTestSuite))
}

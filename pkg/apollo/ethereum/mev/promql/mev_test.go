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
func (t *MevPrometheusTestSuite) TestQueryTopKTokens() {
	t.Require().NotEmpty(t.pc)
	timeNow := time.Now().UTC()
	window := v1.Range{
		Start: timeNow.Add(-time.Minute * 60),
		End:   time.Now().UTC(),
		Step:  time.Minute,
	}
	t.pc.printOn = false
	m, err := t.pc.GetTopTokens(ctx, window)
	t.Require().NoError(err)
	t.Assert().NotEmpty(m)

	for _, val := range m {
		fmt.Println(val.Metric.In, val.Values)
	}
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
	//tfp := web3_client.TokenFeePath{}
	//err = json.Unmarshal(bytes, &tfp)
	//t.Require().NoError(err)
	//fmt.Println("in", tfp.TokenIn.String())
	//fmt.Println("end", tfp.GetEndToken().String())
	//fmt.Println("fee", tfp.Path[0].Fee)

	/*
		{"level":"error","error":"no populated ticks","time":1688261227,"message":"error processing tx"}
		{"level":"error","error":"no populated ticks","path":{"tokenIn":"0x046eee2cc3188071c02bfc1745a6b17c656e3f3d","path":[{"token":"0xdac17f958d2ee523a2206206994597c13d831ec7","fee":3000}]},"time":1688261227,"message":"error getting v3 pricing data"}
	*/
	// TODO, investigate why this is failing
}

func TestMevPrometheusTestSuite(t *testing.T) {
	suite.Run(t, new(MevPrometheusTestSuite))
}

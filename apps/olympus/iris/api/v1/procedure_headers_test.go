package v1_iris

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

type IrisV1TestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *IrisV1TestSuite) TestHeaders() {
	fnRule := iris_programmable_proxy_v1_beta.FanInRuleFirstValidResponse
	ph := ProcedureHeaders{
		XAggOp:               "max",
		XAggKey:              "result",
		XAggKeyValueDataType: "int",
		XAggFilterFanIn:      &fnRule,
	}

	req := &iris_api_requests.ApiProxyRequest{
		Url:             "https://zeus.fyi",
		ExtRoutePath:    "/",
		ServicePlan:     "performance",
		PayloadTypeREST: "POST",
	}
	proc := ph.GetGeneratedProcedure("test", req)
	s.NotNil(proc.OrderedSteps)

	i := 0
	for proc.OrderedSteps.Len() > 0 {
		p := proc.OrderedSteps.PopFront()
		s.NotNil(p)
		ps, ok := p.(iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep)
		s.True(ok)
		if i == 0 {
			s.NotNil(ps.AggregateMap)
		} else {
			s.NotNil(ps.BroadcastInstructions.FanInRules)
		}
		s.Equal("test", ps.BroadcastInstructions.RoutingTable)
		i++
	}
	s.Equal(2, i)
}

func (s *IrisV1TestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	//iris_redis.InitLocalTestProductionRedisIrisCache(ctx)
}

func TestIrisV1TestSuite(t *testing.T) {
	suite.Run(t, new(IrisV1TestSuite))
}

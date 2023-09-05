package v1_iris

import (
	"context"
	"fmt"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_catalog_procedures "github.com/zeus-fyi/olympus/iris/api/v1/procedures"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	"github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

type IrisV1TestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *IrisV1TestSuite) TestEthHeaders() iris_programmable_proxy_v1_beta.IrisRoutingProcedure {
	fnRule := iris_programmable_proxy_v1_beta.FanInRuleFirstValidResponse
	stageTwoPayload := echo.Map{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params": []interface{}{
			"latest",
			true,
		},
		"id": 1,
	}
	ph := ProcedureHeaders{
		XAggOp:                   "max",
		XAggKey:                  "result",
		XAggKeyValueDataType:     "int",
		XAggFilterFanIn:          &fnRule,
		ForwardPayload:           stageTwoPayload,
		StageOneAggregateMapName: iris_catalog_procedures.EthMaxBlockAggReduce,
	}
	payload := iris_catalog_procedures.ProcedureStageOnePayload(iris_catalog_procedures.EthMaxBlockAggReduce)
	req := &iris_api_requests.ApiProxyRequest{
		Url:             "https://fragrant-bitter-lambo.quiknode.pro/7aa59ebefff9128164cac7737bfefeae45b4ee8d/",
		ExtRoutePath:    "/",
		ServicePlan:     "performance",
		PayloadTypeREST: "POST",
		Payload:         payload,
	}
	proc, err := ph.GetGeneratedProcedure("ethereum-mainnet", req)
	s.Nil(err)
	s.NotNil(proc.OrderedSteps)
	s.NotEmpty(payload)
	/*
		needs to test that
		1. max block finder
			a. etl attributes: IrisRoutingResponseETL
				a1. extraction key: eg. result
				a2. extraction key data type: eg. int
				a3. agg op: eg. max
				a4. fan-in rule: eg. returnOnFirstSuccess
		2. needs to update so that it works with json-rpc to replace procedure headers & any method chain
	*/
	return proc
}

func (s *IrisV1TestSuite) TestEthJsonRPC() {
	proc := s.TestEthHeaders()
	s.NotNil(proc)
	groupName := "ethereum-mainnet"
	irisClient := resty_base.GetBaseRestyClient("https://iris.zeus.fyi/v1/router", s.Tc.ProductionLocalTemporalBearerToken)
	irisClient.Header.Set("Content-Type", "application/json")
	irisClient.Header.Set(iris_programmable_proxy.RouteGroupHeader, groupName)
	irisClient.Header.Set(iris_programmable_proxy_v1_beta.LoadBalancingStrategy, iris_programmable_proxy_v1_beta.Adaptive)
	payload := `{
		"jsonrpc": "2.0",
		"procedure": "eth_maxBlockAggReduce",
		"method": "eth_getBlockByNumber",
		"params": ["latest", true],
		"id": 1
	}`
	resp, err := irisClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("/")

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	fmt.Println(resp.String())
}

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
	proc, err := ph.GetGeneratedProcedure("test", req)
	s.Nil(err)
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

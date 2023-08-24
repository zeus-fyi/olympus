package iris_api_requests

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/phf/go-queue/queue"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	iris_operators "github.com/zeus-fyi/zeus/pkg/iris/operators"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

func (s *IrisActivitiesTestSuite) TestBroadcastETL() {
	bc := NewIrisApiRequestsActivities()

	timeOut := time.Second * 3
	rgName := "ethereum-mainnet"
	routes, err := iris_redis.IrisRedisClient.GetBroadcastRoutes(context.Background(), s.Tc.ProductionLocalTemporalOrgID, rgName)
	s.NoError(err)
	s.NotEmpty(routes)

	payload := echo.Map{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  "[]",
		"id":      "1",
	}

	getBlockHeightStep := iris_programmable_proxy_v1_beta.BroadcastInstructions{
		RoutingPath:  "/",
		RestType:     "POST",
		MaxDuration:  timeOut,
		MaxTries:     3,
		RoutingTable: rgName,
		Payload:      payload,
	}

	getBlockHeightProcedure := iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep{
		BroadcastInstructions: getBlockHeightStep,
		TransformSlice: []iris_operators.IrisRoutingResponseETL{
			{
				Source:        "",
				ExtractionKey: "result",
				DataType:      "",
			},
		},
		AggregateMap: map[string]iris_operators.Aggregation{
			"result": {
				Operator:          "max",
				DataType:          "int",
				CurrentMaxInt:     0,
				CurrentMaxFloat64: 0,
			},
		},
	}

	payloadLatestBlock := echo.Map{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{"latest", true},
		"id":      1,
	}

	getBlockStep := iris_programmable_proxy_v1_beta.BroadcastInstructions{
		RoutingPath:  "/",
		RestType:     "POST",
		MaxDuration:  timeOut,
		MaxTries:     3,
		RoutingTable: rgName,
		Payload:      payloadLatestBlock,
		FanInRules: &iris_programmable_proxy_v1_beta.FanInRules{
			Rule: iris_programmable_proxy_v1_beta.FanInRuleFirstValidResponse,
		},
	}

	getBlockProcedure := iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep{
		BroadcastInstructions: getBlockStep,
		TransformSlice: []iris_operators.IrisRoutingResponseETL{
			{
				Source:        "",
				ExtractionKey: "",
				DataType:      "",
			},
		},
	}

	que := queue.New()
	que.PushBack(getBlockHeightProcedure)
	que.PushBack(getBlockProcedure)
	procedure := iris_programmable_proxy_v1_beta.IrisRoutingProcedure{
		Name:         iris_programmable_proxy_v1_beta.MaxBlockAggReduce,
		OrderedSteps: que,
	}

	pr := &ApiProxyRequest{
		Procedure: procedure,
		Routes:    routes,
	}

	fmt.Println(procedure)
	resp, err := bc.BroadcastETLRequest(ctx, pr)
	s.NoError(err)
	s.NotNil(resp)

	s.NotEmpty(resp.Payload)
	fmt.Println(resp.PayloadSizeMeter.Size)
	fmt.Println(resp.Response)
}

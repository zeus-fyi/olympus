package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/phf/go-queue/queue"
	iris_operators "github.com/zeus-fyi/zeus/pkg/iris/operators"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

func (r *IrisRedisTestSuite) TestProcedures() {
	timeOut := time.Second * 3
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
		RoutingTable: "",
		Payload:      payload,
	}

	getBlockHeightProcedure := iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep{
		BroadcastInstructions: getBlockHeightStep,
		TransformSlice: []iris_programmable_proxy_v1_beta.IrisRoutingResponseETL{
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
		RoutingTable: "",
		Payload:      payloadLatestBlock,
		FanInRules: &iris_programmable_proxy_v1_beta.FanInRules{
			Rule: iris_programmable_proxy_v1_beta.FanInRuleFirstValidResponse,
		},
	}

	getBlockProcedure := iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep{
		BroadcastInstructions: getBlockStep,
		TransformSlice: []iris_programmable_proxy_v1_beta.IrisRoutingResponseETL{
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

	err := IrisRedisClient.SetStoredProcedure(context.Background(), 0, procedure)
	r.NoError(err)

	time.Now()
	procedureRetrieved, err := IrisRedisClient.GetStoredProcedure(context.Background(), 0, procedure.Name)
	fmt.Println(time.Since(time.Now()), "time to retrieve")
	r.NoError(err)
	r.Equal(procedure.Name, procedureRetrieved.Name)
	r.Equal(2, procedureRetrieved.OrderedSteps.Len())
}

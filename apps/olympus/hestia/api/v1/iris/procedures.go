package hestia_iris_v1_routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/phf/go-queue/queue"
	iris_operators "github.com/zeus-fyi/zeus/pkg/iris/operators"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

type ProceduresRequest struct {
}

func ProceduresRequestHandler(c echo.Context) error {
	request := new(ProceduresRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetProceduresCatalog(c)
}

func (p *ProceduresRequest) GetProceduresCatalog(c echo.Context) error {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  "[]",
		"id":      "1",
	}
	timeOut := time.Second * 3
	getBlockHeightStep := iris_programmable_proxy_v1_beta.BroadcastInstructions{
		RoutingPath: "/",
		RestType:    "POST",
		MaxDuration: timeOut,
		MaxTries:    3,
		Payload:     payload,
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
	payloadLatestBlock := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{"latest", true},
		"id":      1,
	}
	getBlockStep := iris_programmable_proxy_v1_beta.BroadcastInstructions{
		RoutingPath: "/",
		RestType:    "POST",
		MaxDuration: timeOut,
		MaxTries:    3,
		Payload:     payloadLatestBlock,
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
	resp := []iris_programmable_proxy_v1_beta.IrisRoutingProcedure{{
		Name:         "ethereumMaxBlockAggReduce",
		OrderedSteps: que,
	}, {
		Name:         "near",
		OrderedSteps: que,
	},
	}

	return c.JSON(http.StatusOK, resp)
}

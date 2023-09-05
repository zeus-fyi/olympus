package hestia_iris_v1_routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	iris_operators "github.com/zeus-fyi/zeus/pkg/iris/operators"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

type ProceduresRequest struct {
}

type ProceduresRequestResponseWrapper struct {
	Procedures []ProceduresRequestResponse `json:"procedures"`
}

type ProceduresRequestResponse struct {
	Name         string                                                     `json:"name"`
	Description  string                                                     `json:"description"`
	Protocol     string                                                     `json:"protocol"`
	OrderedSteps []iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep `json:"orderedSteps,omitempty"`
}

func ProceduresRequestHandler(c echo.Context) error {
	request := new(ProceduresRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetProceduresCatalog(c)
}

func ProceduresOnTableRequestHandler(c echo.Context) error {
	request := new(ProceduresRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetProceduresOnTable(c)
}

func (p *ProceduresRequest) GetProceduresOnTable(c echo.Context) error {
	return nil
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
	resp := []ProceduresRequestResponse{{
		Name:         "eth_maxBlockAggReduce",
		Description:  "This procedure will return the latest block number and the block data for the latest block. ",
		Protocol:     "ethereum",
		OrderedSteps: []iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep{getBlockHeightProcedure, getBlockProcedure},
	}, {
		Name:         "near_maxBlockAggReduce",
		Protocol:     "near",
		Description:  "This procedure will return the latest block number and the block data for the latest block. ",
		OrderedSteps: []iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep{getBlockHeightProcedure, getBlockProcedure},
	},
	}

	return c.JSON(http.StatusOK, resp)
}

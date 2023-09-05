package v1_iris

import (
	"encoding/json"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/phf/go-queue/queue"
	"github.com/rs/zerolog/log"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_operators "github.com/zeus-fyi/zeus/pkg/iris/operators"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

func (p *ProxyRequest) ExtractProcedureIfExists(c echo.Context) string {
	procName := c.Request().Header.Get(iris_programmable_proxy_v1_beta.RequestHeaderRoutingProcedureHeader)
	if procName != "" {
		return procName
	}
	procedure := p.Body["procedure"]
	if procedure != nil {
		procNameStr, ok := procedure.(string)
		if ok {
			return procNameStr
		}
	}
	return ""
}

type ProcedureHeaders struct {
	XAggOp               string
	XAggKey              string
	XAggFanIn            *string
	XAggKeyValueDataType string
	XAggComp             string
	XAggCompDataType     string
	XAggFilterPayload    string
	XAggFilterFanIn      *string
}

const (
	xAggOp               = "X-Agg-Op"
	xAggKey              = "X-Agg-Key"
	xAggKeyValueDataType = "X-Agg-Key-Value-Data-Type"
	xAggFanIn            = "X-Agg-Fan-In"
	xAggComp             = "X-Agg-Comp"
	xAggCompDataType     = "X-Agg-Comp-Data-Type"
	xAggFilterPayload    = "X-Agg-Filter-Payload"
	xAggFilterFanIn      = "X-Agg-Filter-Fan-In"
)

func (p *ProcedureHeaders) GetProcedureFromHeaders(c echo.Context, rg string, req *iris_api_requests.ApiProxyRequest) (iris_programmable_proxy_v1_beta.IrisRoutingProcedure, error) {
	if c == nil {
		// for testing
		return p.GetGeneratedProcedure(rg, req)
	}
	ph := GetProcedureFromHeaders(c)
	return ph.GetGeneratedProcedure(rg, req)
}

func (p *ProcedureHeaders) GetGeneratedProcedure(rg string, req *iris_api_requests.ApiProxyRequest) (iris_programmable_proxy_v1_beta.IrisRoutingProcedure, error) {
	if req == nil {
		return iris_programmable_proxy_v1_beta.IrisRoutingProcedure{}, errors.New("req is nil")
	}
	proc := req.Procedure
	if p.XAggKey == "" {
		return proc, errors.New("X-Agg-Key is required")
	}
	switch p.XAggKeyValueDataType {
	case "int", "float64":
	default:
		return proc, errors.New("X-Agg-Key-Value-Data-Type must be int or float64")
	}
	forwardPayload := echo.Map{}
	if p.XAggFilterPayload != "" {
		err := json.Unmarshal([]byte(p.XAggFilterPayload), &forwardPayload)
		if err != nil {
			log.Err(err).Str("payload", p.XAggFilterPayload).Msg("GetGeneratedProcedure: json.Unmarshal")
			return proc, errors.New("X-Agg-Filter-Payload is not valid serialized json")
		}
	}

	var comp *iris_operators.Operation
	if p.XAggComp != "" {
		switch p.XAggCompDataType {
		case "float64", "int":
			comp = &iris_operators.Operation{
				Operator:  iris_operators.Op(p.XAggComp),
				DataTypeX: p.XAggCompDataType,
				DataTypeY: p.XAggCompDataType,
				DataTypeZ: "bool",
			}
		default:
		}
	}
	if comp == nil && p.XAggOp == "" {
		return proc, errors.New("X-Agg-Op is required")
	}
	agg := iris_operators.Aggregation{
		Comparison: comp,
		DataType:   p.XAggKeyValueDataType,
	}
	switch p.XAggOp {
	case "max", "sum":
		agg.Operator = iris_operators.AggOp(p.XAggOp)
	}
	aggValueExtraction := iris_operators.IrisRoutingResponseETL{
		ExtractionKey: p.XAggKey,
		DataType:      p.XAggKeyValueDataType,
	}

	var fanInRules *iris_programmable_proxy_v1_beta.FanInRules
	if p.XAggFanIn != nil {
		fanInRules = &iris_programmable_proxy_v1_beta.FanInRules{
			Rule: iris_programmable_proxy_v1_beta.BroadcastRules(*p.XAggFanIn),
		}
	}
	step := iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep{
		BroadcastInstructions: iris_programmable_proxy_v1_beta.BroadcastInstructions{
			RoutingPath:  req.ExtRoutePath,
			RestType:     req.PayloadTypeREST,
			Payload:      req.Payload,
			MaxDuration:  req.Timeout,
			MaxTries:     req.MaxTries,
			RoutingTable: rg,
			FanInRules:   fanInRules,
		},
		TransformSlice: []iris_operators.IrisRoutingResponseETL{aggValueExtraction},
		AggregateMap: map[string]iris_operators.Aggregation{
			p.XAggKey: agg,
		},
	}
	proc.OrderedSteps = queue.New()
	proc.OrderedSteps.PushBack(step)
	if p.XAggFilterFanIn != nil {
		if *p.XAggFilterFanIn == iris_programmable_proxy_v1_beta.FanInRuleFirstValidResponse {
			fanIn := iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep{
				TransformSlice: []iris_operators.IrisRoutingResponseETL{
					{
						ExtractionKey: "",
						Value:         nil,
					},
				},
				BroadcastInstructions: iris_programmable_proxy_v1_beta.BroadcastInstructions{
					RoutingPath:  req.ExtRoutePath,
					RestType:     req.PayloadTypeREST,
					Payload:      forwardPayload,
					MaxDuration:  req.Timeout,
					MaxTries:     req.MaxTries,
					RoutingTable: rg,
					FanInRules: &iris_programmable_proxy_v1_beta.FanInRules{
						Rule: iris_programmable_proxy_v1_beta.BroadcastRules(*p.XAggFilterFanIn),
					},
				},
			}
			proc.OrderedSteps.PushBack(fanIn)
		}
	}

	req.Procedure = proc
	if req.Procedure.OrderedSteps.Len() == 0 {
		return proc, errors.New("no procedure steps generated")
	}
	return proc, nil
}

func GetProcedureFromHeaders(c echo.Context) *ProcedureHeaders {
	p := &ProcedureHeaders{}
	p.getXAggOp(c)
	p.getXAggKey(c)
	p.getAggFanIn(c)
	p.getXAggKeyValueDataType(c)
	p.getXAggComp(c)
	p.getXAggCompDataType(c)
	p.getXAggFilterPayload(c)
	p.getXAggFilterFanIn(c)
	return p
}

func (p *ProcedureHeaders) getAggFanIn(c echo.Context) {
	x := c.Request().Header.Get(xAggFanIn)
	p.XAggFanIn = &x
}

func (p *ProcedureHeaders) getXAggOp(c echo.Context) {
	x := c.Request().Header.Get(xAggOp)
	p.XAggOp = x
}

func (p *ProcedureHeaders) getXAggKey(c echo.Context) {
	x := c.Request().Header.Get(xAggKey)
	p.XAggKey = x
}

func (p *ProcedureHeaders) getXAggKeyValueDataType(c echo.Context) {
	x := c.Request().Header.Get(xAggKeyValueDataType)
	p.XAggKeyValueDataType = x
}

func (p *ProcedureHeaders) getXAggComp(c echo.Context) {
	x := c.Request().Header.Get(xAggComp)
	p.XAggComp = x
}

func (p *ProcedureHeaders) getXAggCompDataType(c echo.Context) {
	x := c.Request().Header.Get(xAggCompDataType)
	p.XAggCompDataType = x
}

func (p *ProcedureHeaders) getXAggFilterPayload(c echo.Context) {
	x := c.Request().Header.Get(xAggFilterPayload)
	p.XAggFilterPayload = x
}

func (p *ProcedureHeaders) getXAggFilterFanIn(c echo.Context) {
	x := c.Request().Header.Get(xAggFilterFanIn)
	p.XAggFilterFanIn = &x
}

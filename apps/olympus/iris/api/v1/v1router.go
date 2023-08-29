package v1_iris

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phf/go-queue/queue"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	iris_service_plans "github.com/zeus-fyi/olympus/iris/api/v1/service_plans"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_operators "github.com/zeus-fyi/zeus/pkg/iris/operators"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

const QuickNodeMarketPlace = "quickNodeMarketPlace"

func InitV1Routes(e *echo.Echo) {
	eg := e.Group("/v1")
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			// Get headers
			qnTestHeader := c.Request().Header.Get(QuickNodeTestHeader)
			qnIDHeader := c.Request().Header.Get(QuickNodeIDHeader)
			qnEndpointID := c.Request().Header.Get(QuickNodeEndpointID)
			qnChain := c.Request().Header.Get(QuickNodeChain)
			qnNetwork := c.Request().Header.Get(QuickNodeNetwork)

			// Set headers to echo context
			c.Set(QuickNodeTestHeader, qnTestHeader)
			c.Set(QuickNodeIDHeader, qnIDHeader)
			c.Set(QuickNodeEndpointID, qnEndpointID)
			c.Set(QuickNodeChain, qnChain)
			c.Set(QuickNodeNetwork, qnNetwork)
			orgID, plan, err := iris_redis.IrisRedisClient.GetAuthCacheIfExists(ctx, token)
			if err == nil && orgID > 0 && plan != "" {
				c.Set("servicePlan", plan)
				c.Set("orgUser", org_users.NewOrgUserWithID(int(orgID), 0))
				c.Set("bearer", token)
				return true, nil
			} else {
				plan = ""
				orgID = -1
				err = nil
			}
			key := read_keys.NewKeyReader()
			services, err := key.QueryUserAuthedServices(ctx, token)
			if err != nil {
				log.Err(err).Msg("InitV1Routes: QueryUserAuthedServices error")
				return false, err
			}
			c.Set("servicePlans", key.Services)
			if val, ok := key.Services[QuickNodeMarketPlace]; ok {
				plan = val
			} else {
				log.Warn().Str("marketplace", QuickNodeMarketPlace).Msg("InitV1Routes: marketplace not found")
				return false, errors.New("marketplace plan not found")
			}
			ou := org_users.NewOrgUserWithID(key.OrgID, key.GetUserID())
			c.Set("servicePlan", plan)
			c.Set("orgUser", ou)
			c.Set("bearer", token)
			if err == nil && ou.OrgID > 0 && plan != "" {
				go func(oID int, token, plan string) {
					log.Info().Int("orgID", oID).Str("plan", plan).Msg("InitV1Routes: SetAuthCache")
					err = iris_redis.IrisRedisClient.SetAuthCache(context.Background(), oID, token, plan)
					if err != nil {
						log.Err(err).Msg("InitV1Routes: SetAuthCache")
					}
				}(ou.OrgID, token, plan)
			}
			return len(services) > 0, err
		},
	}))

	eg.POST("/router", RpcLoadBalancerPOSTRequestHandler)
	eg.POST("/router/*", wrapHandlerWithCapture(RpcLoadBalancerPOSTRequestHandler))

	eg.GET("/router", RpcLoadBalancerGETRequestHandler)
	eg.GET("/router/*", wrapHandlerWithCapture(RpcLoadBalancerGETRequestHandler))

	eg.PUT("/router", RpcLoadBalancerPUTRequestHandler)
	eg.PUT("/router/*", wrapHandlerWithCapture(RpcLoadBalancerPUTRequestHandler))

	eg.DELETE("/router", RpcLoadBalancerDELETERequestHandler)
	eg.DELETE("/router/*", wrapHandlerWithCapture(RpcLoadBalancerDELETERequestHandler))

	eg.GET(iris_service_plans.PlanUsageDetailsRoute, iris_service_plans.PlanUsageDetailsRequestHandler)
	eg.GET(iris_service_plans.TableMetricsDetailsRoute, iris_service_plans.TableMetricsDetailsRequestHandler)
}

func wrapHandlerWithCapture(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// c.Param("*") will contain the captured path
		capturedPath := c.Param("*")
		fmt.Println(capturedPath)
		c.Set("capturedPath", capturedPath)

		// Then do something with the captured path...
		return handler(c)
	}
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

	var fanInAgg *string
	if p.XAggFanIn != nil {
		fanInAgg = p.XAggFanIn
	}
	step := iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep{
		BroadcastInstructions: iris_programmable_proxy_v1_beta.BroadcastInstructions{
			RoutingPath:  req.ExtRoutePath,
			RestType:     req.PayloadTypeREST,
			Payload:      req.Payload,
			MaxDuration:  req.Timeout,
			MaxTries:     req.MaxTries,
			RoutingTable: rg,
			FanInRules:   &iris_programmable_proxy_v1_beta.FanInRules{Rule: iris_programmable_proxy_v1_beta.BroadcastRules(*fanInAgg)},
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

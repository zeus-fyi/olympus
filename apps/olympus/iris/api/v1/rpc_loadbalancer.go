package v1_iris

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

const (
	QuickNodeTestHeader = "X-QN-TESTING"
	QuickNodeIDHeader   = "x-quicknode-id"
	QuickNodeEndpointID = "x-instance-id"
	QuickNodeChain      = "x-qn-chain"
	QuickNodeNetwork    = "x-qn-network"
	RouteGroupHeader    = "X-Route-Group"
)

const (
	LoadBalancingStrategy    = "X-Load-Balancing-Strategy"
	AnvilSessionLockHeader   = "X-Anvil-Session-Lock-ID"
	EndSessionLockHeader     = "X-End-Session-Lock-ID"
	Adaptive                 = "Adaptive"
	RoundRobin               = "RoundRobin"
	AdaptiveLoadBalancingKey = "X-Adaptive-Metrics-Key"
	EthereumJsonRPC          = "Ethereum"
	QuickNodeJsonRPC         = "QuickNode"
	JsonRpcAdaptiveMetrics   = "JSON-RPC"
)

type ProxyRequest struct {
	Body echo.Map
}

type Response struct {
	Message string `json:"message"`
}

var (
	RpcLoadBalancerGETRequestHandler    = RpcLoadBalancerRequestHandler("GET")
	RpcLoadBalancerPOSTRequestHandler   = RpcLoadBalancerRequestHandler("POST")
	RpcLoadBalancerPUTRequestHandler    = RpcLoadBalancerRequestHandler("PUT")
	RpcLoadBalancerDELETERequestHandler = RpcLoadBalancerRequestHandler("DELETE")
)

func GetDefaultLB(plan string) string {
	switch plan {
	case "enterprise", "standard", "performance", "lite":
		return Adaptive
	default:
		return RoundRobin
	}
}

func RpcLoadBalancerRequestHandler(method string) func(c echo.Context) error {
	return func(c echo.Context) error {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Err(err)
			return err
		}

		payloadSizingMeter := iris_usage_meters.NewPayloadSizeMeter(bodyBytes)
		request := new(ProxyRequest)
		request.Body = echo.Map{}

		if err = json.NewDecoder(payloadSizingMeter).Decode(&request.Body); err != nil {
			log.Err(err).Msgf("RpcLoadBalancerRequestHandler: json.NewDecoder.Decode")
			return err
		}
		anvilHeader := c.Request().Header.Get(AnvilSessionLockHeader)
		if anvilHeader != "" {
			ou := org_users.OrgUser{}
			ouc := c.Get("orgUser")
			if ouc != nil {
				ouser, ok := ouc.(org_users.OrgUser)
				if ok {
					ou = ouser
				} else {
					return c.JSON(http.StatusUnauthorized, Response{Message: "user not found"})
				}
			}
			if ou.OrgID <= 0 {
				return c.JSON(http.StatusUnauthorized, Response{Message: "org not found"})
			}
			plan, ok := c.Get("servicePlan").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, Response{Message: "no service plan found"})
			}
			token, ok := c.Get("bearer").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, Response{Message: "no bearer token found"})
			}
			tempToken, terr := iris_redis.IrisRedisClient.SetInternalAuthCache(context.Background(), ou, token, plan, anvilHeader)
			if terr != nil {
				log.Err(terr).Msg("ProcessRpcLoadBalancerRequest: SetInternalAuthCache")
				return terr
			}

			go func(orgID int, usage *iris_usage_meters.PayloadSizeMeter) {
				err = iris_redis.IrisRedisClient.RecordRequestUsage(context.Background(), orgID, usage)
				if err != nil {
					log.Err(err).Interface("orgID", orgID).Interface("usage", usage).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.IncrementResponseUsageRateMeter")
				}
			}(ou.OrgID, payloadSizingMeter)
			return request.ProcessLockedSessionRoute(c, ou.OrgID, anvilHeader, method, tempToken)
		}

		if payloadSizingMeter.N() <= 0 {
			return request.ProcessRpcLoadBalancerRequest(c, payloadSizingMeter, method)
		}

		lbStrategy := c.Request().Header.Get(LoadBalancingStrategy)
		if lbStrategy == "" {
			lbDefault := c.Get("lbDefault")
			if lbDefault != nil {
				lbStrategyTmp, ok := lbDefault.(string)
				if ok {
					lbStrategy = lbStrategyTmp
				}
			}
		}
		if lbStrategy == Adaptive {
			metric := ""
			adaptiveMetricKeyValue := c.Request().Header.Get(AdaptiveLoadBalancingKey)
			switch adaptiveMetricKeyValue {
			case EthereumJsonRPC, QuickNodeJsonRPC, JsonRpcAdaptiveMetrics:
				metricName := request.Body["method"]
				if metricName != nil {
					metricNameStr, ok := metricName.(string)
					if ok {
						metric = metricNameStr
					}
				}
				c.Set("adaptiveMetricKeyValue", adaptiveMetricKeyValue)
			default:
				metricName := request.Body["method"]
				if metricName != nil {
					metricNameStr, ok := metricName.(string)
					if ok {
						metric = metricNameStr
					}
				}
				adaptiveMetricKeyValue = JsonRpcAdaptiveMetrics
				c.Set("adaptiveMetricKeyValue", adaptiveMetricKeyValue)
			}
			return request.ProcessAdaptiveLoadBalancerRequest(c, payloadSizingMeter, method, metric, adaptiveMetricKeyValue)
		}
		return request.ProcessRpcLoadBalancerRequest(c, payloadSizingMeter, method)
	}
}

func (p *ProxyRequest) ProcessRpcLoadBalancerRequest(c echo.Context, payloadSizingMeter *iris_usage_meters.PayloadSizeMeter, restType string) error {
	procName := p.ExtractProcedureIfExists(c)
	if procName != "" {
		return p.ProcessBroadcastETLRequest(c, payloadSizingMeter, restType, procName, "", "")
	}
	routeGroup := c.Request().Header.Get(RouteGroupHeader)
	if routeGroup == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "routeGroup is required"})
	}
	ou := org_users.OrgUser{}
	ouc := c.Get("orgUser")
	if ouc != nil {
		ouser, ok := ouc.(org_users.OrgUser)
		if ok {
			ou = ouser
		} else {
			return c.JSON(http.StatusBadRequest, Response{Message: "user not found"})
		}
	}
	plan := ""
	if c.Get("servicePlan") != nil {
		sp, ok := c.Get("servicePlan").(string)
		if ok {
			plan = sp
		}
	}
	if plan == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "no service plan found"})
	}
	if payloadSizingMeter == nil {
		payloadSizingMeter = iris_usage_meters.NewPayloadSizeMeter(nil)
	}
	routeInfo, err := iris_redis.IrisRedisClient.CheckRateLimit(context.Background(), ou.OrgID, plan, routeGroup, payloadSizingMeter)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.CheckRateLimit")
		return c.JSON(http.StatusTooManyRequests, Response{Message: err.Error()})
	}
	path := routeInfo.RoutePath
	payloadSizingMeter.Reset()
	headers := make(http.Header)
	for k, v := range c.Request().Header {
		if k == "Authorization" {
			continue // Skip authorization headers
		}
		if k == "X-Route-Group" {
			continue // Skip empty headers
		}
		if k == "X-Load-Balancing-Strategy" {
			continue // Skip empty headers
		}
		if k == "X-Adaptive-Metrics-Key" {
			continue // Skip empty headers
		}
		headers[k] = v // Assuming there's at least one value
	}
	qps := c.QueryParams()
	req := &iris_api_requests.ApiProxyRequest{
		Url:              path,
		ServicePlan:      plan,
		OrgID:            ou.OrgID,
		PayloadTypeREST:  restType,
		RequestHeaders:   headers,
		Referrers:        routeInfo.Referrers,
		Payload:          p.Body,
		QueryParams:      qps,
		IsInternal:       false,
		Timeout:          1 * time.Minute,
		StatusCode:       http.StatusOK, // default
		PayloadSizeMeter: payloadSizingMeter,
	}
	sfx := c.Get("capturedPath")
	if sfx != nil {
		suffix, sok := sfx.(string)
		if sok {
			req.ExtRoutePath = suffix
		}
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	resp, err := rw.ExtLoadBalancerRequest(context.Background(), req)
	if err != nil {
		log.Err(err).Interface("ou", ou).Str("route", path).Msg("ProcessRpcLoadBalancerRequest: rw.ExtLoadBalancerRequest")
		for key, values := range resp.ResponseHeaders {
			for _, value := range values {
				c.Response().Header().Add(key, value)
			}
		}
		for key, values := range resp.FinalResponseHeaders {
			for _, value := range values {
				c.Response().Header().Add(key, value)
			}
		}
		c.Response().Header().Set("X-Response-Latency-Milliseconds", fmt.Sprintf("%d", resp.Latency.Milliseconds()))
		c.Response().Header().Set("X-Response-Received-At-UTC", resp.ReceivedAt.UTC().String())
		c.Response().Header().Set("X-Selected-Route", path)
		return c.JSON(resp.StatusCode, string(resp.RawResponse))
	}
	go func(orgID int, ps *iris_usage_meters.PayloadSizeMeter) {
		err = iris_redis.IrisRedisClient.IncrementResponseUsageRateMeter(context.Background(), ou.OrgID, ps)
		if err != nil {
			log.Err(err).Interface("ou", ou).Str("route", path).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.IncrementResponseUsageRateMeter")
		}
	}(ou.OrgID, payloadSizingMeter)
	for key, values := range resp.ResponseHeaders {
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}
	for key, values := range resp.FinalResponseHeaders {
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}
	c.Response().Header().Set("X-Selected-Route", path)
	c.Response().Header().Set("X-Response-Latency-Milliseconds", fmt.Sprintf("%d", resp.Latency.Milliseconds()))
	c.Response().Header().Set("X-Response-Received-At-UTC", resp.ReceivedAt.UTC().String())
	if resp.Response == nil && resp.RawResponse != nil {
		return c.JSON(resp.StatusCode, string(resp.RawResponse))
	}
	return c.JSON(resp.StatusCode, resp.Response)
}

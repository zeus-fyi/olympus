package v1_iris

import (
	"context"
	"encoding/json"
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

	restTypeGET  = "GET"
	restTypePOST = "POST"
)

type ProxyRequest struct {
	Body echo.Map
}

type Response struct {
	Message string `json:"message"`
}

func RpcLoadBalancerRequestHandler(c echo.Context) error {
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Err(err)
		return err
	}
	payloadSizingMeter := iris_usage_meters.NewPayloadSizeMeter(bodyBytes)
	request := new(ProxyRequest)
	request.Body = echo.Map{}
	if payloadSizingMeter.N() <= 0 {
		return request.DirectProcessRpcLoadBalancerRequest(c, payloadSizingMeter, restTypePOST)
	}
	if err = json.NewDecoder(payloadSizingMeter).Decode(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.DirectProcessRpcLoadBalancerRequest(c, payloadSizingMeter, restTypePOST)
}

func RpcLoadBalancerGETRequestHandler(c echo.Context) error {
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Err(err)
		return err
	}
	payloadSizingMeter := iris_usage_meters.NewPayloadSizeMeter(bodyBytes)
	request := new(ProxyRequest)
	request.Body = echo.Map{}
	if payloadSizingMeter.N() <= 0 {
		return request.DirectProcessRpcLoadBalancerRequest(c, payloadSizingMeter, restTypeGET)
	}
	if err = json.NewDecoder(payloadSizingMeter).Decode(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.DirectProcessRpcLoadBalancerRequest(c, payloadSizingMeter, restTypeGET)
}

func (p *ProxyRequest) DirectProcessRpcLoadBalancerRequest(c echo.Context, payloadSizingMeter *iris_usage_meters.PayloadSizeMeter, restType string) error {
	switch restType {
	case restTypeGET:
		return p.ProcessRpcLoadBalancerRequest(c, payloadSizingMeter, restTypeGET)
	}
	return p.ProcessRpcLoadBalancerRequest(c, payloadSizingMeter, restTypePOST)
}

func (p *ProxyRequest) ProcessRpcLoadBalancerRequest(c echo.Context, payloadSizingMeter *iris_usage_meters.PayloadSizeMeter, restType string) error {
	routeGroup := c.QueryParam("routeGroup")
	if routeGroup == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "routeGroup is required"})
	}
	ou := c.Get("orgUser").(org_users.OrgUser)
	plan := ""
	sp, ok := c.Get("servicePlan").(string)
	if ok {
		plan = sp
	}
	if plan == "" {
		return c.JSON(http.StatusBadRequest, Response{Message: "no service plan found"})
	}
	// todo refactor to fetch auth & plan from redis
	err := iris_redis.IrisRedisClient.CheckRateLimit(context.Background(), ou.OrgID, plan, payloadSizingMeter)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.CheckRateLimit")
		return c.JSON(http.StatusTooManyRequests, Response{Message: err.Error()})
	}

	payloadSizingMeter.Plan = plan
	routeInfo, err := iris_redis.IrisRedisClient.GetNextRoute(context.Background(), ou.OrgID, routeGroup, payloadSizingMeter)
	if err != nil {
		log.Err(err).Interface("ou", ou).Str("routeGroup", routeGroup).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.GetNextRoute")
		errResp := Response{Message: "routeGroup not found"}
		return c.JSON(http.StatusBadRequest, errResp)
	}
	payloadSizingMeter.Reset()
	req := &iris_api_requests.ApiProxyRequest{
		Url:              routeInfo.RoutePath,
		ServicePlan:      plan,
		PayloadTypeREST:  restType,
		Referrers:        routeInfo.Referers,
		Payload:          p.Body,
		Response:         nil,
		IsInternal:       false,
		Timeout:          1 * time.Minute,
		PayloadSizeMeter: payloadSizingMeter,
	}
	rw := iris_api_requests.NewArtemisApiRequestsActivities()
	resp, err := rw.ExtLoadBalancerRequest(context.Background(), req)
	if err != nil {
		log.Err(err).Interface("ou", ou).Str("route", routeInfo.RoutePath).Msg("ProcessRpcLoadBalancerRequest: rw.ExtLoadBalancerRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	go func(orgID int, ps *iris_usage_meters.PayloadSizeMeter) {
		err = iris_redis.IrisRedisClient.IncrementResponseUsageRateMeter(context.Background(), ou.OrgID, ps)
		if err != nil {
			log.Err(err).Interface("ou", ou).Str("route", routeInfo.RoutePath).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.IncrementResponseUsageRateMeter")
		}
	}(ou.OrgID, payloadSizingMeter)
	c.Response().Header().Set("X-Selected-Route", routeInfo.RoutePath)
	return c.JSON(http.StatusOK, resp.Response)
}

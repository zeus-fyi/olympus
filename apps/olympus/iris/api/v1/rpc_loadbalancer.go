package v1_iris

import (
	"bytes"
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
	cr := &iris_usage_meters.PayloadSizeMeter{R: bytes.NewReader(bodyBytes)}
	request := new(ProxyRequest)
	request.Body = echo.Map{}
	if err = json.NewDecoder(cr).Decode(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.ProcessRpcLoadBalancerRequest(c, cr)
}

func (p *ProxyRequest) ProcessRpcLoadBalancerRequest(c echo.Context, payloadSizingMeter *iris_usage_meters.PayloadSizeMeter) error {
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
	routeInfo, err := iris_redis.IrisRedis.GetNextRoute(context.Background(), ou.OrgID, routeGroup)
	if err != nil {
		log.Err(err).Interface("ou", ou).Str("routeGroup", routeGroup).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.GetNextRoute")
		errResp := Response{Message: "routeGroup not found"}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	req := &iris_api_requests.ApiProxyRequest{
		Url:              routeInfo.RoutePath,
		ServicePlan:      plan,
		Referrers:        routeInfo.Referers,
		Payload:          p.Body,
		Response:         nil,
		IsInternal:       false,
		Timeout:          1 * time.Minute,
		PayloadSizeMeter: payloadSizingMeter,
	}
	rw := iris_api_requests.NewArtemisApiRequestsActivities()
	resp, err := rw.ExtLoadBalancerRequest(c.Request().Context(), req)
	if err != nil {
		log.Err(err).Interface("ou", ou).Str("route", routeInfo.RoutePath).Msg("ProcessRpcLoadBalancerRequest: rw.ExtLoadBalancerRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resp.Response)
}

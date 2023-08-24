package v1_iris

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func (p *ProxyRequest) ProcessBroadcastETLRequest(c echo.Context, payloadSizingMeter *iris_usage_meters.PayloadSizeMeter, restType, metricName string) error {
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
	ri, err := iris_redis.IrisRedisClient.CheckRateLimit(context.Background(), ou.OrgID, plan, routeGroup, payloadSizingMeter)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("ProcessAdaptiveLoadBalancerRequest: iris_redis.CheckRateLimit")
		return c.JSON(http.StatusTooManyRequests, Response{Message: err.Error()})
	}
	//fmt.Println(ri.RoutePath, "routeRoundRobin")
	routes, err := iris_redis.IrisRedisClient.GetBroadcastRoutes(context.Background(), ou.OrgID, routeGroup)
	if err != nil {
		log.Err(err).Interface("ou", ou).Str("routeGroup", routeGroup).Msg("ProcessAdaptiveLoadBalancerRequest: iris_round_robin.GetNextRoute")
		errResp := Response{Message: "routeGroup not found"}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	//fmt.Println(tableStats.MemberRankScoreOut.Member, "routeAdaptive")
	//fmt.Println(tableStats, "tableStats")
	payloadSizingMeter.Plan = plan
	extPath := ""
	sfx := c.Get("capturedPath")
	if sfx != nil {
		suffix, sok := sfx.(string)
		if sok {
			extPath = suffix
		}
	}
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
		Routes:           routes,
		ExtRoutePath:     extPath,
		ServicePlan:      plan,
		PayloadTypeREST:  restType,
		RequestHeaders:   headers,
		Referrers:        ri.Referers,
		Payload:          p.Body,
		QueryParams:      qps,
		IsInternal:       false,
		Timeout:          1 * time.Minute,
		StatusCode:       http.StatusOK, // default
		PayloadSizeMeter: payloadSizingMeter,
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	var sendRawResponse bool
	now := time.Now()
	resp, err := rw.BroadcastETLRequest(context.Background(), req)
	if err != nil {
		log.Err(err).Interface("ou", ou).Interface("routes", routes).Msg("ProcessRpcLoadBalancerRequest: rw.ExtLoadBalancerRequest")
		err = nil
		sendRawResponse = true
	}
	go func(orgID int, usage *iris_usage_meters.PayloadSizeMeter) {
		err = iris_redis.IrisRedisClient.RecordRequestUsage(context.Background(), orgID, usage)
		if err != nil {
			log.Err(err).Interface("ou", ou).Interface("usage", usage).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.IncrementResponseUsageRateMeter")
		}
	}(ou.OrgID, resp.PayloadSizeMeter)
	for key, values := range resp.ResponseHeaders {
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}
	sinceMs := time.Since(now).Milliseconds()
	c.Response().Header().Set("X-Procedure-Latency-Milliseconds", fmt.Sprintf("%d", sinceMs))
	c.Response().Header().Set("X-Response-Received-At-UTC", resp.ReceivedAt.UTC().String())
	if (resp.Response == nil && resp.RawResponse != nil) || sendRawResponse {
		return c.JSON(resp.StatusCode, string(resp.RawResponse))
	}
	return c.JSON(resp.StatusCode, resp.Response)
}

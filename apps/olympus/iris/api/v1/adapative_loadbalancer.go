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

func (p *ProxyRequest) ProcessAdaptiveLoadBalancerRequest(c echo.Context, payloadSizingMeter *iris_usage_meters.PayloadSizeMeter, restType, metricName, adaptiveKeyName string) error {
	procName := p.ExtractProcedureIfExists(c)
	if procName != "" {
		return p.ProcessBroadcastETLRequest(c, payloadSizingMeter, restType, procName, metricName, adaptiveKeyName)
	}
	route := ""
	routeHost := c.Get("proxy")
	if routeHost != nil {
		tmp, ok := routeHost.(string)
		if ok {
			route = tmp
		}
	}
	routeGroup := c.Request().Header.Get(RouteGroupHeader)
	if routeGroup == "" && route == "" {
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
		} else {
			plan = "free"
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
	tableStats, err := iris_redis.IrisRedisClient.GetNextAdaptiveRoute(context.Background(), ou.OrgID, routeGroup, metricName, ri, payloadSizingMeter)
	if err != nil {
		log.Err(err).Interface("ou", ou).Str("routeGroup", routeGroup).Msg("ProcessAdaptiveLoadBalancerRequest: iris_round_robin.GetNextRoute")
		errResp := Response{Message: "routeGroup not found"}
		return c.JSON(http.StatusBadRequest, errResp)
	}

	//fmt.Println(tableStats.MemberRankScoreOut.Member, "routeAdaptive")
	//fmt.Println(tableStats, "tableStats")
	payloadSizingMeter.Plan = plan
	if tableStats.MemberRankScoreOut.Member == nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	path, ok := tableStats.MemberRankScoreOut.Member.(string)
	if !ok && route == "" {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if route != "" {
		path = route
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

	secretNameRefApi := fmt.Sprintf("api-%s", routeGroup)
	qps := c.QueryParams()
	if len(route) > 0 {
		qps.Del("url")
	}
	req := &iris_api_requests.ApiProxyRequest{
		Url:              path,
		ServicePlan:      plan,
		OrgID:            ou.OrgID,
		UserID:           ou.UserID,
		PayloadTypeREST:  restType,
		RequestHeaders:   headers,
		Referrers:        ri.Referrers,
		Payload:          p.Body,
		QueryParams:      qps,
		IsInternal:       false,
		Timeout:          2 * time.Minute,
		StatusCode:       http.StatusOK, // default
		PayloadSizeMeter: payloadSizingMeter,
		SecretNameRef:    secretNameRefApi,
	}
	sfx := c.Get("capturedPath")
	if sfx != nil {
		suffix, sok := sfx.(string)
		if sok {
			req.ExtRoutePath = suffix
		}
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	var sendRawResponse bool
	resp, err := rw.ExtLoadBalancerRequest(context.Background(), req)
	if err != nil {
		log.Err(err).Interface("ou", ou).Str("route", path).Msg("ProcessRpcLoadBalancerRequest: rw.ExtLoadBalancerRequest")
		err = nil
		sendRawResponse = true
	}
	if resp.StatusCode >= 400 {
		tableStats.LatencyScaleFactor = tableStats.ErrorScaleFactor
	}
	tableStats.Meter = resp.PayloadSizeMeter
	tableStats.LatencyMilliseconds = resp.Latency.Milliseconds()
	//fmt.Println(tableStats, "tableStats")
	go func(orgID int, tbl *iris_redis.StatTable) {
		err = iris_redis.IrisRedisClient.SetLatestAdaptiveEndpointPriorityScoreAndUpdateRateUsage(context.Background(), tbl)
		if err != nil {
			log.Err(err).Interface("ou", ou).Str("route", path).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.IncrementResponseUsageRateMeter")
		}
	}(ou.OrgID, tableStats)
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
	if (resp.Response == nil && resp.RawResponse != nil) || sendRawResponse {
		return c.JSON(resp.StatusCode, string(resp.RawResponse))
	}
	return c.JSON(resp.StatusCode, resp.Response)
}

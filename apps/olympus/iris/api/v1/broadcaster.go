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
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

func (p *ProxyRequest) ProcessBroadcastETLRequest(c echo.Context, payloadSizingMeter *iris_usage_meters.PayloadSizeMeter, restType, procName string) error {
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
	if procName != iris_programmable_proxy_v1_beta.MaxBlockAggReduce {
		procName = ""
	}
	if len(procName) <= 0 {
		procName = fmt.Sprintf("%d", ou.OrgID)
	}
	proc, routes, err := iris_redis.IrisRedisClient.CheckRateLimitBroadcast(context.Background(), ou.OrgID, procName, plan, routeGroup, payloadSizingMeter)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("ProcessAdaptiveLoadBalancerRequest: iris_redis.CheckRateLimit")
		return c.JSON(http.StatusTooManyRequests, nil)
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

	// TODO: add max timeout param
	req := &iris_api_requests.ApiProxyRequest{
		Procedure:        proc,
		Routes:           routes,
		ExtRoutePath:     extPath,
		ServicePlan:      plan,
		PayloadTypeREST:  restType,
		RequestHeaders:   headers,
		Payload:          p.Body,
		QueryParams:      qps,
		IsInternal:       false,
		Timeout:          4 * time.Second,
		StatusCode:       http.StatusOK, // default
		PayloadSizeMeter: payloadSizingMeter,
	}
	orgIDStr := fmt.Sprintf("%d", ou.OrgID)
	if req.Procedure.Name == orgIDStr {
		procHeaders := ProcedureHeaders{}
		_, perr := procHeaders.GetProcedureFromHeaders(c, routeGroup, req)
		if perr != nil {
			return c.JSON(http.StatusBadRequest, Response{Message: perr.Error()})
		}
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
	for key, values := range resp.FinalResponseHeaders {
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}
	c.Response().Header().Set("X-Procedure-Latency-Milliseconds", fmt.Sprintf("%d", sinceMs))
	//c.Response().Header().Set("X-Response-Received-At-UTC", resp.ReceivedAt.UTC().String())
	if (resp.Response == nil && resp.RawResponse != nil) || sendRawResponse {
		return c.JSON(resp.StatusCode, string(resp.RawResponse))
	}
	return c.JSON(resp.StatusCode, resp.Response)
}

package v1_iris

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func ProcessLockedSessionsHandler(c echo.Context) error {
	request := new(ProxyRequest)
	request.Body = echo.Map{}
	if err := c.Bind(&request.Body); err != nil {
		log.Err(err).Msg("proxy_anvil.ProcessLockedSessionsHandler: c.Bind")
		return err
	}
	anvilHeader := c.Request().Header.Get(AnvilSessionLockHeader)
	ou := org_users.OrgUser{}
	ouc := c.Get("orgUser")
	if ouc != nil {
		ouser, ok := ouc.(org_users.OrgUser)
		if ok {
			ou = ouser
		} else {
			log.Warn().Interface("ou", ouser).Msg("proxy_anvil.ProcessLockedSessionsHandler: orgUser not found")
			return c.JSON(http.StatusUnauthorized, Response{Message: "user not found"})
		}
	}
	return request.ProcessLockedSessionRoute(c, ou.OrgID, anvilHeader)
}

/*
TODO: implement this in redis

func (a *AnvilProxy) GetSessionLockedRoute(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "Zeus-Test" {
		return "http://anvil.eeb335ad-78da-458f-9cfb-9928514d65d0.svc.cluster.local:8545", nil
	}
*/

/*
1. needs to have a registry of all anvil routes
  - this should be auto inserted when new anvil services are deployed

2a. needs to have a registry of all locked sessions
2b. needs to have a registry of all locked routes
3. needs to be able to dynamically add/remove anvil services
*/
const (
	anvilServerlessRoutesTableName = "anvil"
)

func GetSessionLockedRoute(ctx context.Context, orgID int, sessionID, tableRoute string) (string, error) {
	if sessionID == "Zeus-Test" {
		return "http://anvil.eeb335ad-78da-458f-9cfb-9928514d65d0.svc.cluster.local:8545", nil
	}
	route, err := iris_redis.IrisRedisClient.GetNextServerlessRoute(context.Background(), orgID, sessionID, tableRoute)
	if err != nil {
		log.Err(err).Msg("proxy_anvil.SessionLocker.GetNextServerlessRoute")
		return route, err
	}
	return route, err
}

func (p *ProxyRequest) ProcessLockedSessionRoute(c echo.Context, orgID int, sessionID string) error {
	endLockedSessionLease := c.Request().Header.Get("End-Session-Lock-ID")
	if endLockedSessionLease == sessionID {
		// todo remove hardcoded table name
		return p.ProcessEndSessionLock(c, orgID, endLockedSessionLease, anvilServerlessRoutesTableName)
	}
	routeURL, err := GetSessionLockedRoute(c.Request().Context(), orgID, sessionID, anvilServerlessRoutesTableName) // TODO remove hardcoded table name
	if err != nil {
		log.Err(err).Msg("proxy_anvil.SessionLocker.GetSessionLockedRoute")
		return c.JSON(http.StatusInternalServerError, err)
	}
	req := &iris_api_requests.ApiProxyRequest{
		Url:        routeURL,
		Payload:    p.Body,
		IsInternal: false,
		Timeout:    60 * time.Second,
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	resp, err := rw.ExtToAnvilInternalSimForkRequest(c.Request().Context(), req)
	if err != nil {
		log.Err(err).Str("routeURL", routeURL).Interface("resp", resp).Msg("rw.InternalSvcRelayRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	go func(orgID int, usage *iris_usage_meters.PayloadSizeMeter) {
		if usage == nil {
			log.Warn().Msg("ProcessLockedSessionRoute: usage is nil")
			return
		}
		err = iris_redis.IrisRedisClient.RecordRequestUsage(context.Background(), orgID, usage)
		if err != nil {
			log.Err(err).Interface("orgID", orgID).Interface("usage", usage).Msg("ProcessRpcLoadBalancerRequest: iris_round_robin.IncrementResponseUsageRateMeter")
		}
	}(orgID, resp.PayloadSizeMeter)
	return c.JSON(http.StatusOK, resp.Response)
}

func (p *ProxyRequest) ProcessEndSessionLock(c echo.Context, orgID int, sessionID, serverlessRoutesTable string) error {
	err := iris_redis.IrisRedisClient.ReleaseServerlessRoute(context.Background(), orgID, sessionID, serverlessRoutesTable)
	if err != nil {
		log.Err(err).Msg("proxy_anvil.SessionLocker.ReleaseServerlessRoute")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func (p *ProxyRequest) Process(c echo.Context, r *iris_api_requests.ApiProxyRequest) error {
	if r == nil {
		return c.JSON(http.StatusBadRequest, errors.New("request is nil"))
	}
	resp, err := iris_api_requests.IrisProxyWorker.ExecuteIrisProxyWorkflow(c.Request().Context(), r)
	if err != nil {
		log.Err(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	if resp == nil {
		if r != nil {
			log.Warn().Str("route", r.Url).Msg("rw.InternalSvcRelayRequest")
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, resp.Response)
}

package v1_iris

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	proxy_anvil "github.com/zeus-fyi/olympus/pkg/iris/proxy/anvil"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
)

func ProcessLockedSessionsHandler(c echo.Context) error {
	request := new(ProxyRequest)
	request.Body = echo.Map{}
	if err := c.Bind(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.ProcessLockedSessionRoute(c)
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

func MockGetSessionLockedRoute(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "Zeus-Test" {
		return "http://anvil.eeb335ad-78da-458f-9cfb-9928514d65d0.svc.cluster.local:8545", nil
	}
	return "", errors.New("not implemented")
}

func (p *ProxyRequest) ProcessLockedSessionRoute(c echo.Context) error {
	sessionID := c.Request().Header.Get("Session-Lock-ID")
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, errors.New("Session-Lock-ID header is required"))
	}
	endLockedSessionLease := c.Request().Header.Get("End-Session-Lock-ID")
	if endLockedSessionLease == sessionID {
		return nil //p.ProcessEndSessionLock(c, endLockedSessionLease)
	}
	routeInfo, err := MockGetSessionLockedRoute(c.Request().Context(), sessionID)
	if err != nil {
		log.Err(err).Msg("proxy_anvil.SessionLocker.GetSessionLockedRoute")
		return c.JSON(http.StatusInternalServerError, err)
	}
	req := &iris_api_requests.ApiProxyRequest{
		Url:        routeInfo,
		Payload:    p.Body,
		IsInternal: false,
		Timeout:    30 * time.Second,
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	resp, err := rw.InternalSvcRelayRequest(c.Request().Context(), req)
	if err != nil {
		log.Err(err).Str("route", routeInfo).Msg("rw.InternalSvcRelayRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resp.Response)
}

func (p *ProxyRequest) ProcessEndSessionLock(c echo.Context, sessionID string) error {
	proxy_anvil.SessionLocker.RemoveSessionLockedRoute(c.Request().Context(), sessionID)
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

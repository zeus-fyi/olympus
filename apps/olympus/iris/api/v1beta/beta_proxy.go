package v1Beta_iris

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	proxy_anvil "github.com/zeus-fyi/olympus/pkg/iris/proxy/anvil"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
)

type BetaProxyRequest struct {
	Body echo.Map
}

func InternalBetaProxyRequestHandler(c echo.Context) error {
	request := new(BetaProxyRequest)
	request.Body = echo.Map{}
	if err := c.Bind(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.ProcessInternalHardhat(c, true)
}

func (p *BetaProxyRequest) ProcessInternalHardhat(c echo.Context, isInternal bool) error {
	sessionID := c.Request().Header.Get("Session-Lock-ID")
	wfExecutor := c.Request().Header.Get("Durable-Execution-ID")
	if sessionID == "" {
		return c.JSON(http.StatusBadRequest, errors.New("Session-Lock-ID header is required"))
	}
	endLockedSessionLease := c.Request().Header.Get("End-Session-Lock-ID")
	if endLockedSessionLease == sessionID {
		return p.ProcessEndSessionLock(c, endLockedSessionLease)
	}
	routeInfo, err := proxy_anvil.SessionLocker.GetSessionLockedRoute(c.Request().Context(), sessionID)
	if err != nil {
		log.Err(err).Msg("proxy_anvil.SessionLocker.GetSessionLockedRoute")
		return c.JSON(http.StatusInternalServerError, err)
	}
	req := &iris_api_requests.ApiProxyRequest{
		Url:        routeInfo,
		Payload:    p.Body,
		IsInternal: isInternal,
		Timeout:    30 * time.Second,
	}
	if wfExecutor != "" {
		return p.Process(c, req)
	}
	rw := iris_api_requests.NewIrisApiRequestsActivities()
	resp, err := rw.InternalSvcRelayRequest(c.Request().Context(), req)
	if err != nil {
		log.Err(err).Str("route", routeInfo).Msg("rw.InternalSvcRelayRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resp.Response)
}

func (p *BetaProxyRequest) ProcessEndSessionLock(c echo.Context, sessionID string) error {
	proxy_anvil.SessionLocker.RemoveSessionLockedRoute(c.Request().Context(), sessionID)
	return c.JSON(http.StatusOK, nil)
}

func (p *BetaProxyRequest) Process(c echo.Context, r *iris_api_requests.ApiProxyRequest) error {
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

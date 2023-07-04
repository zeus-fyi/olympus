package v1_iris

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_api_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/api_requests"
	proxy_anvil "github.com/zeus-fyi/olympus/pkg/iris/proxy/anvil"
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
	req := &artemis_api_requests.ApiProxyRequest{
		Url:        routeInfo,
		Payload:    p.Body,
		IsInternal: isInternal,
		Timeout:    1 * time.Minute,
	}
	if wfExecutor != "" {
		return p.Process(c, req)
	}
	rw := artemis_api_requests.NewArtemisApiRequestsActivities()
	resp, err := rw.InternalSvcRelayRequest(c.Request().Context(), req)
	if err != nil {
		log.Err(err).Msg("rw.InternalSvcRelayRequest")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resp.Response)
}

func (p *BetaProxyRequest) ProcessEndSessionLock(c echo.Context, sessionID string) error {
	proxy_anvil.SessionLocker.RemoveSessionLockedRoute(sessionID)
	return c.JSON(http.StatusOK, nil)
}

func (p *BetaProxyRequest) Process(c echo.Context, r *artemis_api_requests.ApiProxyRequest) error {
	if r == nil {
		return c.JSON(http.StatusBadRequest, errors.New("request is nil"))
	}
	resp, err := artemis_api_requests.ArtemisProxyWorker.ExecuteArtemisInternalSvcApiProxyWorkflow(c.Request().Context(), r)
	if err != nil {
		log.Err(err)
		return err
	}
	if resp == nil {
		log.Warn().Msg("resp == nil")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, resp.Response)
}

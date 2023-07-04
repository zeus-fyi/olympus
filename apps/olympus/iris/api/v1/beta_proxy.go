package v1_iris

import (
	"errors"
	"math/rand"
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
	return request.ProcessInternalHardhat(c, true, 0)
}

func (p *BetaProxyRequest) ProcessInternalHardhat(c echo.Context, isInternal bool, tryNumber int) error {
	maxRetries := 5 // Specify max retries
	if tryNumber > maxRetries {
		return c.JSON(http.StatusInternalServerError, errors.New("max retries exceeded"))
	}
	relayTo := c.Request().Header.Get("Session-Lock-ID")
	if relayTo == "" {
		return c.JSON(http.StatusBadRequest, errors.New("Session-Lock-ID header is required"))
	}
	endLockedSessionLease := c.Request().Header.Get("End-Session-Lock-ID")
	if endLockedSessionLease == relayTo {
		return p.ProcessEndSessionLock(c, endLockedSessionLease)
	}
	r, err := proxy_anvil.SessionLocker.GetSessionLockedRoute(relayTo)
	if err != nil {
		log.Err(err).Interface("relayDestination", relayTo).Msgf("proxy_anvil.SessionLocker.GetSessionLockedRoute %e", err)
		minDuration := 25 * time.Millisecond
		maxDuration := 100 * time.Millisecond
		jitter := time.Duration(1) * (time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration)
		time.Sleep(jitter)
		return p.ProcessInternalHardhat(c, isInternal, tryNumber+1)
	}
	req := &artemis_api_requests.ApiProxyRequest{
		Url:        r.Route,
		Payload:    p.Body,
		IsInternal: isInternal,
		Timeout:    1 * time.Minute,
	}
	wfExecutor := c.Request().Header.Get("Durable-Execution-ID")
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

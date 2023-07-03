package v1_iris

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_api_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/api_requests"
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
	relayTo := c.Request().Header.Get("Session-Lock-ID")
	if relayTo == "" {
		return c.JSON(http.StatusBadRequest, errors.New("Session-Lock-ID header is required"))
	}

	rw := artemis_api_requests.NewArtemisApiRequestsActivities()
	_, err := rw.RelayRequest(c.Request().Context(), &artemis_api_requests.ApiProxyRequest{
		Url:        "https://hardhat.zeus.fyi/",
		IsInternal: isInternal,
		Timeout:    1 * time.Minute,
	})
	if err != nil {
		log.Err(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}

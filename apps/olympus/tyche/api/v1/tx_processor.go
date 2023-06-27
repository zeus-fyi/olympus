package v1_tyche

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
)

type TxProcessingRequest struct {
	Body echo.Map
}

func TxProcessingRequestHandler(c echo.Context) error {
	request := new(TxProcessingRequest)
	request.Body = echo.Map{}
	if err := c.Bind(&request.Body); err != nil {
		log.Err(err)
		return err
	}
	return request.ProcessTx(c)
}

func (t *TxProcessingRequest) ProcessTx(c echo.Context) error {
	a := artemis_realtime_trading.ActiveTrading{}
	a.ProcessTx(c.Request().Context(), nil)
	return nil
}

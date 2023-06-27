package v1_tyche

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
)

type TxProcessingRequest struct {
	Txs []*types.Transaction
}

func TxProcessingRequestHandler(c echo.Context) error {
	request := new(TxProcessingRequest)
	if err := c.Bind(&request); err != nil {
		log.Err(err)
		return err
	}
	return request.ProcessTx(c)
}

func (t *TxProcessingRequest) ProcessTx(c echo.Context) error {
	a := artemis_realtime_trading.ActiveTrading{}

	// should process in parallel
	a.ProcessTx(c.Request().Context(), nil)
	return nil
}

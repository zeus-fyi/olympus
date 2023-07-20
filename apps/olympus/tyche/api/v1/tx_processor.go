package v1_tyche

import (
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	artemis_trade_executor "github.com/zeus-fyi/olympus/pkg/artemis/trading/executor"
	tyche_metrics "github.com/zeus-fyi/olympus/tyche/metrics"
)

const (
	txProcessorRoute = "/v1/mev/mempool/tx"
)

type TxProcessingRequest struct {
	SeenAt time.Time            `json:"seenAt"`
	Txs    []*types.Transaction `json:"txs"`
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
	ctx := c.Request().Context()
	w3c := artemis_trade_executor.ActiveTraderW3c
	for _, tx := range t.Txs {
		werr := artemis_realtime_trading.IngestTx(ctx, w3c, tx, &tyche_metrics.TradeMetrics)
		if werr.Err != nil && werr.Code != 200 {
			//log.Err(werr.Err).Msg("error processing tx")
			return c.JSON(http.StatusPreconditionFailed, werr.Err)
		}
	}
	return c.JSON(http.StatusOK, "ok")
}

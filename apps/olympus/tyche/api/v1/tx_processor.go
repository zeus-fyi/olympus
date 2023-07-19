package v1_tyche

import (
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_trade_executor "github.com/zeus-fyi/olympus/pkg/artemis/trading/executor"
	tyche_metrics "github.com/zeus-fyi/olympus/tyche/metrics"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
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
	at := artemis_trade_executor.ActiveTrader
	for _, tx := range t.Txs {
		switch int(tx.ChainId().Int64()) {
		case hestia_req_types.EthereumGoerliProtocolNetworkID:
			at = artemis_trade_executor.ActiveGoerliTrader
		case hestia_req_types.EthereumMainnetProtocolNetworkID:
			at = artemis_trade_executor.ActiveTrader
		case hestia_req_types.EthereumEphemeryProtocolNetworkID:
			log.Info().Msgf("tx chain id %s not supported or not supplied, defaulting to mainnet", tx.ChainId().String())
		default:
			log.Info().Msgf("tx chain id %s not supported or not supplied, defaulting to mainnet", tx.ChainId().String())
		}
		werr := at.IngestTx(ctx, tx, &tyche_metrics.TradeMetrics)
		if werr.Err != nil && werr.Code != 200 {
			//log.Err(werr.Err).Msg("error processing tx")
			return c.JSON(http.StatusPreconditionFailed, werr.Err)
		}
	}
	return c.JSON(http.StatusOK, "ok")
}

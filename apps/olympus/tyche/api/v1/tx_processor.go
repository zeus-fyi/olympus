package v1_tyche

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	tyche_metrics "github.com/zeus-fyi/olympus/tyche/metrics"
)

const (
	txProcessorRoute = "/v1/mev/mempool/tx"
)

type TxProcessingRequest struct {
	Txs []*types.Transaction `json:"txs"`
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
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)
	uni := web3_client.InitUniswapClient(ctx, wc)
	a := artemis_realtime_trading.NewActiveTradingModule(&uni, tyche_metrics.TradeMetrics)
	// should process in parallel

	for _, tx := range t.Txs {
		a.IngestTx(ctx, tx)
	}
	return nil
}

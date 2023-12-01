package v1_tyche

import (
	"context"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
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
	//w3cTrader := artemis_trade_executor.ActiveTraderW3c
	for _, tx := range t.Txs {
		//go func(tx *types.Transaction, w3c web3_client.Web3Client) {
		//	werr := artemis_realtime_trading.IngestTx(context.Background(), w3c, tx, &tyche_metrics.TradeMetrics)
		//	if werr.Err != nil && werr.Code != 200 {
		//		log.Err(werr.Err).Msg("error processing tx")
		//	}
		//}(tx, w3cTrader)
		go func(tx *types.Transaction) {
			if tx == nil {
				return
			}
			b, err := tx.MarshalBinary()
			if err != nil {
				log.Err(err).Msg("error marshalling tx")
				return
			}
			err = iris_redis.IrisRedisClient.CreateOrAddToStream(context.Background(), iris_redis.EthMempoolStreamName, map[string]interface{}{
				tx.Hash().Hex(): b,
			})
			if err != nil {
				log.Err(err).Msg("error adding to redis mempool stream")
				return
			}
		}(tx)
	}
	return c.JSON(http.StatusOK, "ok")
}

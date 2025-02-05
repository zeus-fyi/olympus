package artemis_realtime_trading

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func SaveMempoolTx(ctx context.Context, tfSlice []web3_client.TradeExecutionFlow, m *metrics_trading.TradingMetrics) error {
	for _, tradeFlow := range tfSlice {
		if !tradeFlow.AreAllTradesValid() {
			continue
		}
		tradeJSON, err := tradeFlow.ConvertToJSONType()
		if err != nil {
			return err
		}
		btf, ber := json.Marshal(tradeJSON)
		if ber != nil {
			log.Err(ber).Msg("failed to marshal tx flow")
			return ber
		}
		fromStr := ""
		chainId := artemis_eth_units.NewBigInt(hestia_req_types.EthereumMainnetProtocolNetworkID)
		if tradeFlow.Tx.ChainId() != nil {
			chainId = tradeFlow.Tx.ChainId()
		}
		sender := types.LatestSignerForChainID(chainId)
		from, ferr := sender.Sender(tradeFlow.Tx)
		if ferr != nil {
			log.Err(ferr).Msg("failed to get sender")
			return ferr
		} else {
			fromStr = from.String()
		}
		txStr, terr := json.Marshal(tradeJSON.Tx)
		if terr != nil {
			log.Err(terr).Msg("failed to marshal tx")
			return terr
		}

		nonce := tradeFlow.Tx.Nonce()
		txMempool := artemis_autogen_bases.EthMempoolMevTx{
			ProtocolNetworkID: hestia_req_types.EthereumMainnetProtocolNetworkID,
			Tx:                string(txStr),
			TxFlowPrediction:  string(btf),
			TxHash:            tradeFlow.Tx.Hash().String(),
			Nonce:             int(nonce),
			From:              fromStr,
			To:                tradeFlow.Tx.To().String(),
			BlockNumber:       int(tradeFlow.CurrentBlockNumber.Uint64()),
		}
		err = artemis_mev_models.InsertMempoolTx(ctx, txMempool)
		if err != nil {
			log.Err(err).Msg("failed to insert mempool tx")
			return err
		}
		if m != nil {
			m.StageProgressionMetrics.CountSavedMempoolTx(float64(1))
		}
	}
	return nil
}

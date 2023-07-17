package artemis_realtime_trading

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

// SaveMempoolTx Also filters WETH denominated txs for active trading consideration
func SaveMempoolTx(ctx context.Context, bn uint64, tfSlice []web3_client.TradeExecutionFlowJSON, m *metrics_trading.TradingMetrics) error {
	var liveTradingSlice []web3_client.TradeExecutionFlowJSON
	for _, tradeFlow := range tfSlice {
		if tradeFlow.SandwichPrediction.ExpectedProfit == "0" || tradeFlow.SandwichPrediction.ExpectedProfit == "1" {
			continue
		}
		if tradeFlow.UserTrade.AmountIn == "" {
			continue
		}
		if tradeFlow.SandwichTrade.AmountOut == "" {
			continue
		}
		if tradeFlow.UserTrade.AmountInAddr.String() == artemis_trading_constants.WETH9ContractAddressAccount.String() {
			liveTradingSlice = append(liveTradingSlice, tradeFlow)
		}
		tradeFlow.CurrentBlockNumber = new(big.Int).SetUint64(bn)
		btf, ber := json.Marshal(tradeFlow)
		if ber != nil {
			log.Err(ber).Msg("failed to marshal tx flow")
			return ber
		}
		baseTx, zerr := tradeFlow.Tx.ConvertToTx()
		if zerr != nil {
			log.Err(zerr).Msg("dat: EntryTxFilter, ConvertToTx")
			return zerr
		}
		fromStr := ""
		sender := types.LatestSignerForChainID(baseTx.ChainId())
		from, ferr := sender.Sender(baseTx)
		if ferr != nil {
			log.Err(ferr).Msg("failed to get sender")
			return ferr
		} else {
			fromStr = from.String()
		}
		txStr, terr := json.Marshal(tradeFlow.Tx)
		if terr != nil {
			log.Err(terr).Msg("failed to marshal tx")
			return terr
		}
		txMempool := artemis_autogen_bases.EthMempoolMevTx{
			ProtocolNetworkID: hestia_req_types.EthereumMainnetProtocolNetworkID,
			Tx:                string(txStr),
			TxFlowPrediction:  string(btf),
			TxHash:            tradeFlow.Tx.Hash,
			Nonce:             int(baseTx.Nonce()),
			From:              fromStr,
			To:                tradeFlow.Tx.To,
			BlockNumber:       int(bn),
		}
		err := artemis_mev_models.InsertMempoolTx(ctx, txMempool)
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

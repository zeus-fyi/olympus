package artemis_trade_debugger

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (t *TradeDebugger) lookupMevTx(ctx context.Context, txHash string) ([]HistoricalAnalysisDebug, error) {
	mevTxs, merr := artemis_mev_models.SelectEthMevTxAnalysisByTxHash(ctx, txHash)
	if merr != nil {
		return nil, merr
	}
	historicalAnalysisDebugs := make([]HistoricalAnalysisDebug, len(mevTxs))
	for i, mevTx := range mevTxs {
		tfPrediction, err := web3_client.UnmarshalTradeExecutionFlow(mevTx.TxFlowPrediction)
		if err != nil {
			return nil, err
		}
		historicalAnalysisDebugs[i] = HistoricalAnalysisDebug{
			HistoricalAnalysis: mevTx,
			TradePrediction:    tfPrediction.ConvertToBigIntType(),
		}
		switch mevTx.TradeMethod {
		case artemis_trading_constants.V2SwapExactIn:
			trade := historicalAnalysisDebugs[i].TradePrediction.Trade.JSONV2SwapExactInParams.ConvertToBigIntType()
			historicalAnalysisDebugs[i].TradeParams = trade
		}
	}
	return historicalAnalysisDebugs, merr
}

func (t *TradeDebugger) getTxFromHash(ctx context.Context, txHash string) (*types.Transaction, error) {
	hash := common.HexToHash(txHash)
	tx, _, err := t.UniswapClient.Web3Client.GetTxByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (t *TradeDebugger) getRxFromHash(ctx context.Context, txHash string) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	rx, err := t.UniswapClient.Web3Client.GetTxReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}
	return rx, nil
}

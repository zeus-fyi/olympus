package artemis_trade_debugger

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (t *TradeDebugger) lookupMevMempoolTx(ctx context.Context, txHash string) (HistoricalAnalysisDebug, error) {
	mevMempoolTx, err := artemis_mev_models.SelectEthMevMempoolTxByTxHash(ctx, txHash)
	if err != nil {
		return HistoricalAnalysisDebug{}, err
	}
	if len(mevMempoolTx) == 0 {
		return HistoricalAnalysisDebug{}, errors.New("no mev tx found")
	}
	//historicalAnalysisDebugs := make([]HistoricalAnalysisDebug, len(mevMempoolTx))
	for _, mevTx := range mevMempoolTx {
		tfPrediction, terr := web3_client.UnmarshalTradeExecutionFlow(mevTx.TxFlowPrediction)
		if terr != nil {
			return HistoricalAnalysisDebug{}, terr
		}
		tp, terr := tfPrediction.ConvertToBigIntType()
		if terr != nil {
			return HistoricalAnalysisDebug{}, terr
		}
		dbgTx := HistoricalAnalysisDebug{
			HistoricalAnalysis: mevTx,
			TradePrediction:    tp,
		}
		return dbgTx, err
		//historicalAnalysisDebugs[i] = dbgTx
	}

	return HistoricalAnalysisDebug{}, errors.New("no mev tx found")
}

func GetMevMempoolTxTradeFlow(ctx context.Context, txHash string) (*web3_client.TradeExecutionFlowJSON, error) {
	mevMempoolTx, err := artemis_mev_models.SelectEthMevMempoolTxByTxHash(ctx, txHash)
	if err != nil {
		return nil, err
	}
	for _, mevTx := range mevMempoolTx {
		tfPrediction, terr := web3_client.UnmarshalTradeExecutionFlow(mevTx.TxFlowPrediction)
		if terr == nil {
			return &tfPrediction, terr
		}
	}
	return nil, err
}

func (t *TradeDebugger) lookupMevTx(ctx context.Context, txHash string) (HistoricalAnalysisDebug, error) {
	mevTxs, merr := artemis_mev_models.SelectEthMevTxAnalysisByTxHash(ctx, txHash)
	if merr != nil {
		return HistoricalAnalysisDebug{}, merr
	}
	if len(mevTxs) == 0 {
		return HistoricalAnalysisDebug{}, errors.New("no mev tx found")

	}
	historicalAnalysisDebugs := make([]HistoricalAnalysisDebug, len(mevTxs))
	for i, mevTx := range mevTxs {
		tfPrediction, err := web3_client.UnmarshalTradeExecutionFlow(mevTx.TxFlowPrediction)
		if err != nil {
			return HistoricalAnalysisDebug{}, err
		}
		tp, err := tfPrediction.ConvertToBigIntType()
		if err != nil {
			return HistoricalAnalysisDebug{}, err
		}
		wrapper := HistoricalAnalysisDebug{
			HistoricalAnalysis: mevTx,
			TradePrediction:    tp,
		}
		historicalAnalysisDebugs[i] = wrapper
		switch mevTx.TradeMethod {
		case artemis_trading_constants.V2SwapExactIn:
			trade, terr := historicalAnalysisDebugs[i].TradePrediction.Trade.JSONV2SwapExactInParams.ConvertToBigIntType()
			if terr != nil {
				return HistoricalAnalysisDebug{}, terr
			}
			historicalAnalysisDebugs[i].TradeParams = trade
		}
		if tfPrediction.Tx.Hash == txHash {
			return wrapper, nil
		}
	}
	return HistoricalAnalysisDebug{}, merr
}

func (t *TradeDebugger) getTxFromHash(ctx context.Context, txHash string) (*types.Transaction, error) {
	hash := common.HexToHash(txHash)
	tx, _, err := t.dat.SimW3c().GetTxByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (t *TradeDebugger) getRxFromHash(ctx context.Context, txHash string) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	rx, _, err := t.dat.SimW3c().GetTxReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}
	return rx, nil
}

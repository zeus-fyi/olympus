package artemis_trade_debugger

import (
	"context"
	"fmt"
)

func (t *TradeDebugger) getMevTx(ctx context.Context, txHash string, fromMempoolTx bool) (HistoricalAnalysisDebug, error) {
	if fromMempoolTx {
		return t.lookupMevMempoolTx(ctx, txHash)
	}
	return t.lookupMevTx(ctx, txHash)
}
func (t *TradeDebugger) Replay(ctx context.Context, txHash string, fromMempoolTx bool) error {
	mevTx, err := t.getMevTx(ctx, txHash, fromMempoolTx)
	if err != nil {
		return err
	}
	tf := mevTx.TradePrediction
	err = t.ResetAndSetupPreconditions(ctx, tf)
	if err != nil {
		return err
	}
	fmt.Println("ANALYZING tx: ", tf.Tx.Hash().String(), "at block: ", mevTx.GetBlockNumber())
	amountInAddr := tf.FrontRunTrade.AmountInAddr
	amountIn := tf.FrontRunTrade.AmountIn
	err = t.analyzeToken(ctx, amountInAddr, amountIn)
	if err != nil {
		return err
	}
	amountOutAddr := tf.FrontRunTrade.AmountOutAddr
	amountOut := tf.FrontRunTrade.AmountOut
	err = t.analyzeToken(ctx, amountOutAddr, amountOut)
	if err != nil {
		return err
	}
	_, err = t.UniswapClient.FrontRunTradeGetAmountsOut(&tf)
	if err != nil {
		err = t.analyzeDrift(ctx, tf.FrontRunTrade)
		return err
	}
	ac := t.ActiveTrading.GetAuxClient()
	ur, err := ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.FrontRunTrade)
	if err != nil {
		return err
	}
	if ur == nil {
		return fmt.Errorf("ur is nil")
	}
	err = t.UniswapClient.InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.FrontRunTrade)
	if err != nil {
		return err
	}
	//err = t.analyzeUserTrade(ctx, &tf)
	//if err != nil {
	//	return err
	//}
	_, err = t.UniswapClient.ExecTradeByMethod(&tf)
	if err != nil {
		return err
	}
	ur, err = ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.SandwichTrade)
	if err != nil {
		return err
	}
	if ur == nil {
		return fmt.Errorf("ur is nil")
	}
	err = t.UniswapClient.InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.SandwichTrade)
	if err != nil {
		return err
	}
	err = t.UniswapClient.VerifyTradeResults(&tf)
	if err != nil {
		return err
	}
	return nil
}

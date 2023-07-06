package artemis_trade_debugger

import (
	"context"
	"fmt"
)

func (t *TradeDebugger) Replay(ctx context.Context, txHash string) error {
	mevTx, err := t.lookupMevTx(ctx, txHash)
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
	err = t.UniswapClient.ExecTradeV2SwapFromTokenToToken(ctx, &tf.FrontRunTrade)
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
	err = t.UniswapClient.ExecTradeV2SwapFromTokenToToken(ctx, &tf.SandwichTrade)
	if err != nil {
		return err
	}
	err = t.UniswapClient.VerifyTradeResults(&tf)
	if err != nil {
		return err
	}
	return nil
}

package artemis_trade_debugger

import (
	"context"
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
	amountInAddr := tf.FrontRunTrade.AmountInAddr
	amountIn := tf.FrontRunTrade.AmountIn
	err = t.analyzeToken(ctx, amountInAddr, amountIn)
	if err != nil {
		return err
	}
	_, err = t.UniswapClient.ExecFrontRunTradeStepTokenTransfer(&tf)
	if err != nil {
		return err
	}
	return nil
}

package artemis_trade_debugger

import (
	"context"
	"fmt"

	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
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
	_, err = t.dat.GetSimUniswapClient().FrontRunTradeGetAmountsOut(&tf)
	if err != nil {
		err = t.analyzeDrift(ctx, tf.FrontRunTrade)
		return err
	}
	ac := t.dat.GetSimAuxClient()
	tf.FrontRunTrade.AmountOut = tf.FrontRunTrade.SimulatedAmountOut //  new(big.Int).SetInt64(0)
	ur, err := ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.FrontRunTrade)
	if err != nil {
		return err
	}
	if ur == nil {
		return fmt.Errorf("ur is nil")
	}
	err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.FrontRunTrade)
	if err != nil {
		return err
	}
	_, err = t.dat.GetSimUniswapClient().ExecTradeByMethod(&tf)
	if err != nil {
		return err
	}
	tf.SandwichTrade.AmountIn = tf.FrontRunTrade.SimulatedAmountOut
	_, err = t.dat.GetSimUniswapClient().SandwichTradeGetAmountsOut(&tf)
	if err != nil {
		err = t.analyzeDrift(ctx, tf.FrontRunTrade)
		return err
	}
	startBal, err := ac.CheckAuxERC20BalanceFromAddr(ctx, tf.SandwichTrade.AmountOutAddr.String())
	if err != nil {
		return err
	}

	tf.SandwichTrade.AmountOut = tf.SandwichTrade.SimulatedAmountOut
	ur, err = ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.SandwichTrade)
	if err != nil {
		return err
	}
	if ur == nil {
		return fmt.Errorf("ur is nil")
	}
	err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.SandwichTrade)
	if err != nil {
		return err
	}
	endBal, err := ac.CheckAuxERC20BalanceFromAddr(ctx, tf.SandwichTrade.AmountOutAddr.String())
	if err != nil {
		return err
	}
	fmt.Println("profit", artemis_eth_units.SubBigInt(endBal, startBal))
	//err = t.dat.GetSimUniswapClient().VerifyTradeResults(&tf)
	//if err != nil {
	//	return err
	//}
	return nil
}

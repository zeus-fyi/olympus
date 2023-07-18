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
	ac := t.dat.GetSimAuxClient()
	ur, err := ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.FrontRunTrade)
	if err != nil {
		return err
	}
	err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.FrontRunTrade)
	if err != nil {
		err = t.FindSlippage(ctx, &tf.FrontRunTrade)
		if err != nil {
			return err
		}
	}
	_, err = t.dat.GetSimUniswapClient().ExecTradeByMethod(&tf)
	if err != nil {
		return err
	}
	startBal, err := ac.CheckAuxERC20BalanceFromAddr(ctx, tf.SandwichTrade.AmountOutAddr.String())
	if err != nil {
		return err
	}
	tf.SandwichTrade.AmountIn = tf.FrontRunTrade.AmountOut
	ur, err = ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.SandwichTrade)
	if err != nil {
		return err
	}
	if ur == nil {
		return fmt.Errorf("ur is nil")
	}
	err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.SandwichTrade)
	if err != nil {
		err = t.FindSlippage(ctx, &tf.SandwichTrade)
		if err != nil {
			return err
		}
	}
	endBal, err := ac.CheckAuxERC20BalanceFromAddr(ctx, tf.SandwichTrade.AmountOutAddr.String())
	if err != nil {
		return err
	}
	fmt.Println("profitToken", tf.SandwichTrade.AmountOutAddr.String())
	fmt.Println("expectedProfit", tf.SandwichTrade.AmountOut.String())
	fmt.Println("actualProfit", artemis_eth_units.SubBigInt(endBal, startBal))
	return nil
}

/*
	frontRunTokenInAddr := tf.FrontRunTrade.AmountInAddr.String()
	if info, ok := artemis_trading_cache.TokenMap[frontRunTokenInAddr]; ok {
		den := info.TransferTaxDenominator
		num := info.TransferTaxNumerator
		if den != nil && num != nil {
			fmt.Println("token: ", frontRunTokenInAddr, "tradingTax: num: ", *num, "den: ", *den)
		} else {
			fmt.Println("token not found in cache")
		}
		tf.FrontRunTrade.AmountOut = artemis_eth_units.ApplyTransferTax(tf.FrontRunTrade.AmountOut, *num, *den)
	}
	frontRunTokenOutAddr := tf.FrontRunTrade.AmountOutAddr.String()
	if info, ok := artemis_trading_cache.TokenMap[frontRunTokenOutAddr]; ok {
		den := info.TransferTaxDenominator
		num := info.TransferTaxNumerator
		if den != nil && num != nil {
			fmt.Println("token: ", frontRunTokenOutAddr, "tradingTax: num: ", *num, "den: ", *den)
		} else {
			fmt.Println("token not found in cache")
		}
		tf.FrontRunTrade.AmountOut = artemis_eth_units.ApplyTransferTax(tf.FrontRunTrade.AmountOut, *num, *den)
	}
	tf.SandwichTrade.AmountIn = tf.FrontRunTrade.AmountOut

*/

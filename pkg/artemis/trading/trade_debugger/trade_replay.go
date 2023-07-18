package artemis_trade_debugger

import (
	"context"
	"fmt"

	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
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
	fmt.Println("FRONT RUN TRADE: ", tf.FrontRunTrade.AmountInAddr.String(), " -> ", tf.FrontRunTrade.AmountOutAddr.String())
	ac := t.dat.GetSimAuxClient()
	n, d := GetMaxTransferTax(tf)
	amountOutStartFrontRun := tf.FrontRunTrade.AmountOut
	amountOutStartSandwich := tf.SandwichTrade.AmountOut

	adjAmountOut := artemis_eth_units.ApplyTransferTax(amountOutStartFrontRun, n, d)
	tf.FrontRunTrade.AmountOut = adjAmountOut
	ur, err := ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.FrontRunTrade)
	if err != nil {
		return err
	}
	err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.FrontRunTrade)
	if err != nil {
		tf.FrontRunTrade.AmountOut = amountOutStartFrontRun
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
	adjAmountOut = artemis_eth_units.ApplyTransferTax(amountOutStartSandwich, n+30, d)
	tf.SandwichTrade.AmountOut = adjAmountOut
	ur, err = ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.SandwichTrade)
	if err != nil {
		return err
	}
	if ur == nil {
		return fmt.Errorf("ur is nil")
	}
	err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.SandwichTrade)
	if err != nil {
		tf.SandwichTrade.AmountOut = amountOutStartSandwich
		err = t.FindSlippage(ctx, &tf.SandwichTrade)
		if err != nil {
			return err
		}
		return err
	}
	endBal, err := ac.CheckAuxERC20BalanceFromAddr(ctx, tf.SandwichTrade.AmountOutAddr.String())
	if err != nil {
		return err
	}

	profitToken := tf.SandwichTrade.AmountOutAddr.String()
	fmt.Println("profitToken", tf.SandwichTrade.AmountOutAddr.String())
	fmt.Println("expectedProfit", tf.SandwichTrade.AmountOut.String())
	expProfit := artemis_eth_units.SubBigInt(endBal, startBal)
	fmt.Println("expProfit", expProfit)

	err = tf.GetAggregateGasUsage(ctx, ac.U.Web3Client)
	if err != nil {
		return err
	}
	totalGasCost := tf.SandwichTrade.TotalGasCost + tf.FrontRunTrade.TotalGasCost
	fmt.Println("totalGasCost", totalGasCost)

	if profitToken == artemis_trading_constants.WETH9ContractAddress {
		expProfit = artemis_eth_units.SubUint64FBigInt(expProfit, totalGasCost)
	}
	err = artemis_mev_models.UpdateEthMevTxAnalysis(ctx, txHash, expProfit.String(), fmt.Sprintf("%d", totalGasCost), "success")
	if err != nil {
		return err
	}
	return nil
}

func GetMaxTransferTax(tf web3_client.TradeExecutionFlow) (int, int) {
	tokenOne := tf.UserTrade.AmountInAddr.String()
	tokenTwo := tf.UserTrade.AmountOutAddr.String()
	maxNum, maxDen := 0, 1
	if info, ok := artemis_trading_cache.TokenMap[tokenOne]; ok {
		den := info.TransferTaxDenominator
		num := info.TransferTaxNumerator
		if den != nil && num != nil {
			fmt.Println("token: ", tokenOne, "transferTax: num: ", *num, "den: ", *den)

			if *num > maxNum {
				maxNum = *num
				maxDen = *den
			}
		} else {
			fmt.Println("token not found in cache")
		}
	}
	if info, ok := artemis_trading_cache.TokenMap[tokenTwo]; ok {
		den := info.TransferTaxDenominator
		num := info.TransferTaxNumerator
		if den != nil && num != nil {
			fmt.Println("token: ", tokenTwo, "tradingTax: num: ", *num, "den: ", *den)
			if *num > maxNum {
				maxNum = *num
				maxDen = *den
			}
		} else {
			fmt.Println("token not found in cache")
		}
	}
	fmt.Println("maxNum: ", maxNum, "maxDen: ", maxDen)
	return maxNum, maxDen
}

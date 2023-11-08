package artemis_trade_debugger

import (
	"context"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
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
	fmt.Println("FRONT_RUN AmountIn:", tf.FrontRunTrade.AmountIn.String(), "AmountInAddr", tf.FrontRunTrade.AmountInAddr.String())
	fmt.Println("BACK_RUN AmountOut", tf.SandwichTrade.AmountOut.String(), "AmountOutAddr", tf.SandwichTrade.AmountOutAddr.String())
	fmt.Println("TRADE_METHOD:", tf.Trade.TradeMethod, "TradeMethod")
	fmt.Println("EXPECTED_PROFIT", tf.SandwichPrediction.ExpectedProfit, "AmountOutAddr", tf.SandwichTrade.AmountOutAddr.String())

	err = t.ResetAndSetupPreconditions(context.Background(), tf)
	if err != nil {
		return err
	}
	fmt.Println("ANALYZING tx: ", tf.Tx.Hash().String(), "at block: ", mevTx.GetBlockNumber())
	fmt.Println("FRONT_RUN TRADE: ", tf.FrontRunTrade.AmountInAddr.String(), " -> ", tf.FrontRunTrade.AmountOutAddr.String())
	ac := t.dat.GetSimAuxClient()
	n, d := GetMaxTransferTax(tf)
	amountOutStartFrontRun := tf.FrontRunTrade.AmountOut
	amountOutStartSandwich := tf.SandwichTrade.AmountOut

	adjAmountOut := artemis_eth_units.ApplyTransferTax(amountOutStartFrontRun, n, d)
	tf.FrontRunTrade.AmountOut = adjAmountOut
	w3c := ac.U.Web3Client
	ur, _, err := artemis_trading_auxiliary.GenerateTradeV2SwapFromTokenToToken(ctx, w3c, nil, &tf.FrontRunTrade)
	if err != nil {
		return err
	}
	err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.FrontRunTrade)
	if err != nil {
		tf.FrontRunTrade.AmountOut = amountOutStartFrontRun
		err = t.FindSlippage(ctx, w3c, &tf.FrontRunTrade)
		if err != nil {
			log.Err(err).Str("txHash", txHash).Msg("FRONT_RUN: error finding slippage")
			return err
		}
	}
	_, err = t.dat.GetSimUniswapClient().ExecTradeByMethod(&tf)
	if err != nil {
		log.Err(err).Str("txHash", txHash).Msg("USER_TRADE: error executing user trade")
		return err
	}
	startBal, err := ac.CheckAuxERC20BalanceFromAddr(ctx, tf.SandwichTrade.AmountOutAddr.String())
	if err != nil {
		log.Err(err).Str("txHash", txHash).Msg("error checking balance")
		return err
	}
	tf.SandwichTrade.AmountIn = tf.FrontRunTrade.AmountOut

	if d == 1000 {
		n += 30
	}
	adjAmountOut = artemis_eth_units.ApplyTransferTax(amountOutStartSandwich, n, d)
	tf.SandwichTrade.AmountOut = adjAmountOut
	ur, _, err = artemis_trading_auxiliary.GenerateTradeV2SwapFromTokenToToken(ctx, w3c, nil, &tf.SandwichTrade)
	if err != nil || ur == nil {
		if err == nil {
			err = fmt.Errorf("ur is nil")
		}
		log.Err(err).Str("txHash", txHash).Msg("error finding slippage")
		return err
	}
	err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.SandwichTrade)
	if err != nil {
		tf.SandwichTrade.AmountOut = amountOutStartSandwich
		err = t.FindSlippage(ctx, w3c, &tf.SandwichTrade)
		if err != nil {
			log.Err(err).Str("txHash", txHash).Msg("SANDWICH_TRADE: error finding slippage")
			return err
		}
	}
	endBal, err := ac.CheckAuxERC20BalanceFromAddr(ctx, tf.SandwichTrade.AmountOutAddr.String())
	if err != nil {
		log.Err(err).Str("txHash", txHash).Msg("error checking balance")
		return err
	}
	fmt.Println("DONE ANALYZING tx: ", tf.Tx.Hash().String(), "at block: ", mevTx.GetBlockNumber())
	profitToken := tf.SandwichTrade.AmountOutAddr.String()
	fmt.Println("profitToken", profitToken)
	fmt.Println("sandwichTfAmountOut", tf.SandwichTrade.AmountOut.String())

	expProfit := artemis_eth_units.SubBigInt(endBal, startBal)
	fmt.Println("expProfitAmountOutBalanceChange", expProfit)

	if artemis_eth_units.IsXLessThanY(tf.SandwichTrade.AmountOut, expProfit) {
		expProfit = tf.SandwichTrade.AmountOut
	}
	log.Info().Str("txHash", txHash).Str("expProfitWorstCase", expProfit.String()).Msg("expected profit")
	err = tf.GetAggregateGasUsage(ctx, ac.U.Web3Client)
	if err != nil {
		log.Err(err).Str("txHash", txHash).Msg("error getting gas usage")
		return err
	}
	totalGasCost := tf.SandwichTrade.TotalGasCost + tf.FrontRunTrade.TotalGasCost
	fmt.Println("totalGasCost", totalGasCost)

	if profitToken == artemis_trading_constants.WETH9ContractAddress {
		expProfit = artemis_eth_units.SubUint64FBigInt(expProfit, totalGasCost)
	}

	if t.insertNewTxs {
		rx, rerr := t.getRxFromHash(ctx, txHash)
		if rerr != nil {
			return rerr
		}
		pair := ""
		if tf.InitialPair.PairContractAddr != "" {
			pair = tf.InitialPair.PairContractAddr
		} else {
			if tf.InitialPairV3.PoolAddress != "" {
				pair = tf.InitialPairV3.PoolAddress
			}
		}
		tradeAnalysis := web3_client.TradeAnalysisReport{
			TxHash:             txHash,
			TradeMethod:        tf.Trade.TradeMethod,
			ArtemisBlockNumber: int(tf.CurrentBlockNumber.Int64()),
			RxBlockNumber:      int(rx.BlockNumber.Int64()),
			PairAddress:        pair,
			GasReport: web3_client.GasReport{
				TotalGasUsed:         fmt.Sprintf("%d", totalGasCost),
				FrontRunGasUsed:      fmt.Sprintf("%d", tf.FrontRunTrade.TotalGasCost),
				SandwichTradeGasUsed: fmt.Sprintf("%d", tf.SandwichTrade.TotalGasCost),
			},
			TradeFailureReport: web3_client.TradeFailureReport{
				EndReason: "success",
				EndStage:  "",
			},
			SimulationResults: web3_client.SimulationResults{
				AmountInAddr:            tf.FrontRunTrade.AmountInAddr.String(),
				AmountIn:                tf.FrontRunTrade.AmountIn.String(),
				AmountOutAddr:           tf.SandwichTrade.AmountOutAddr.String(),
				AmountOut:               new(big.Int).Sub(tf.SandwichTrade.AmountOut, tf.FrontRunTrade.AmountIn).String(),
				ExpectedProfitAmountOut: expProfit.String(),
			},
		}
		tradeAnalysis.PrintResults()
		err = tradeAnalysis.SaveResultsInDb(ctx)
		if err != nil {
			return err
		}
	} else {
		err = artemis_mev_models.UpdateEthMevTxAnalysis(ctx, txHash, expProfit.String(), fmt.Sprintf("%d", totalGasCost), "success")
		if err != nil {
			return err
		}
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

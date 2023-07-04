package web3_client

import (
	"errors"
	"fmt"
	"math/big"
)

/*
Trade Method: swapExactTokensForETH
Artemis Block Number: 17390714
Rx Block Number: 17390715
End Reason: unable to overwrite balance
End Stage: executing front run balance setup

Trade Method: swapExactTokensForETH
Artemis Block Number: 17390721
Rx Block Number: 17390723
End Reason: unable to overwrite balance
End Stage: executing front run balance setup

verifying full sandwich trade
Trade Method: swapExactTokensForETH
Artemis Block Number: 17390766
Rx Block Number: 17390767
End Reason: user trade amount out mismatch
*/

func (u *UniswapClient) VerifyTradeResults(tf *TradeExecutionFlow) error {
	if u.DebugPrint {
		fmt.Println("verifying full sandwich trade")
	}

	/*
		func (t *artemis_trading_types.TradeOutcome) GetGasUsageForAllTxs(ctx context.Context, w Web3Client) error {
			for _, tx := range t.OrderedTxs {
				txInfo, err := w.GetTxLifecycleStats(ctx, accounts.HexToHash(tx.Hex()))
				if err != nil {
					return err
				}
				t.TotalGasCost += txInfo.GasUsed
			}
			return nil
		}

	*/
	switch tf.Trade.TradeMethod {
	case swapTokensForExactETH:

	case swapExactTokensForETH:
		/*
				Trade Method: swapExactTokensForETH
				Artemis Block Number: 17390700
				Rx Block Number: 17390701
				End Reason: user trade amount out mismatch

			Trade Method: swapExactTokensForETH
			Artemis Block Number: 17390791
			Rx Block Number: 17390792
			End Reason: unable to overwrite balance
			End Stage: executing front run balance setup
		*/
		//gasAdjustedBalance := new(big.Int).Sub(tf.UserTrade.AmountOut, new(big.Int).SetUint64(tf.UserTrade.TotalGasCost))
		//difference := new(big.Int).Sub(tf.UserTrade.PostTradeEthBalance, tf.UserTrade.PreTradeEthBalance)
		//if difference.String() != gasAdjustedBalance.String() {
		//	return errors.New("user trade amount out mismatch")
		//}
	case swapExactTokensForTokens:
	case swapTokensForExactTokens:
	case swapExactETHForTokens:
	case swapETHForExactTokens:

	default:
		//return errors.New("invalid trade method")
	}

	frontRunGasCost := new(big.Int).SetUint64(tf.FrontRunTrade.TotalGasCost)
	u.TradeAnalysisReport.GasReport.FrontRunGasUsed = frontRunGasCost.String()
	fmt.Println("frontRunGasCost", frontRunGasCost.String())

	sandwichRunGasCost := new(big.Int).SetUint64(tf.SandwichTrade.TotalGasCost)
	u.TradeAnalysisReport.GasReport.SandwichTradeGasUsed = sandwichRunGasCost.String()
	fmt.Println("sandwichRunGasCost", sandwichRunGasCost.String())

	totalSandwichTradeGasCost := new(big.Int).Add(frontRunGasCost, sandwichRunGasCost)
	fmt.Println("total gas cost", totalSandwichTradeGasCost.String())
	u.TotalGasUsed = totalSandwichTradeGasCost.String()

	gasFreeProfit := new(big.Int).Sub(tf.SandwichTrade.AmountOut, tf.FrontRunTrade.AmountIn)
	fmt.Println("gas free profit", gasFreeProfit.String(), "profitToken", tf.SandwichTrade.AmountOutAddr.String())
	expMinusActualProfit := new(big.Int).Sub(tf.SandwichPrediction.ExpectedProfit, gasFreeProfit)
	if expMinusActualProfit.String() != "0" {
		return errors.New("expected minus actual profit mismatch")
	}
	if tf.SandwichTrade.AmountOutAddr.String() == WETH9ContractAddress {
		realizedProfit := new(big.Int).Sub(gasFreeProfit, totalSandwichTradeGasCost)
		fmt.Println("realized profit", realizedProfit.String())
	}

	u.EndReason = "success"
	return nil
}

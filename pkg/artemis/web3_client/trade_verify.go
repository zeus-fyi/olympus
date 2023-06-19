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
		gasAdjustedBalance := new(big.Int).Sub(tf.UserTrade.AmountOut, new(big.Int).SetUint64(tf.UserTrade.TotalGasCost))
		difference := new(big.Int).Sub(tf.UserTrade.PostTradeEthBalance, tf.UserTrade.PreTradeEthBalance)
		if difference.String() != gasAdjustedBalance.String() {
			return errors.New("user trade amount out mismatch")
		}
	case swapExactTokensForTokens:
	case swapTokensForExactTokens:
	case swapExactETHForTokens:
	case swapETHForExactTokens:

	default:
		return errors.New("invalid trade method")
	}

	frontRunGasCost := new(big.Int).SetUint64(tf.FrontRunTrade.TotalGasCost)
	u.TradeAnalysisReport.GasReport.FrontRunGasUsed = frontRunGasCost.String()

	sandwichRunGasCost := new(big.Int).SetUint64(tf.SandwichTrade.TotalGasCost)
	u.TradeAnalysisReport.GasReport.SandwichTradeGasUsed = frontRunGasCost.String()

	totalSandwichTradeGasCost := new(big.Int).Add(frontRunGasCost, sandwichRunGasCost)
	fmt.Println("total gas cost", totalSandwichTradeGasCost.String())
	u.TotalGasUsed = totalSandwichTradeGasCost.String()

	endingTokenBalance, err := u.Web3Client.ReadERC20TokenBalance(ctx, tf.SandwichTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		return err
	}
	fmt.Println("profit currency", tf.SandwichTrade.AmountOutAddr.String())
	u.AmountOutAddr = tf.SandwichTrade.AmountOutAddr.String()
	fmt.Println("starting amount", tf.FrontRunTrade.AmountIn.String())
	fmt.Println("ending amount", tf.SandwichTrade.AmountOut.String())
	profitTokenBalance := new(big.Int).Sub(endingTokenBalance, tf.FrontRunTrade.AmountIn)
	fmt.Println("profitTokenBalance", profitTokenBalance.String())
	fmt.Println("sandwichCalculatedProfit", tf.SandwichPrediction.ExpectedProfit.String())
	u.SimulationResults.AmountOut = profitTokenBalance.String()
	u.SimulationResults.ExpectedProfitAmountOut = tf.SandwichPrediction.ExpectedProfit.String()
	if profitTokenBalance.String() != tf.SandwichPrediction.ExpectedProfit.String() {
		fmt.Println("profit token balance mismatch", "profitTokenBalance", profitTokenBalance.String(), "expectedProfit", tf.SandwichPrediction.ExpectedProfit.String())
		diff := new(big.Int).Sub(profitTokenBalance, tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("diff", "diff", diff.String())
		diff2 := new(big.Int).Sub(tf.SandwichPrediction.ExpectedProfit, profitTokenBalance)
		fmt.Println("diff2", "diff2", diff2.String())

		return errors.New("profit token balance mismatch")
	}
	u.EndReason = "success"

	return nil
}

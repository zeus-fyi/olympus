package web3_client

import (
	"errors"
	"fmt"
	"math/big"
)

func (u *UniswapV2Client) VerifyTradeResults(tf *TradeExecutionFlowInBigInt) error {
	if u.DebugPrint {
		fmt.Println("verifying full sandwich trade")
	}

	switch tf.Trade.TradeMethod {
	case swapTokensForExactETH:
	case swapExactTokensForETH:
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

	if tf.SandwichTrade.AmountOutAddr.String() != WETH9ContractAddress {
		fmt.Println("profit currency", tf.SandwichTrade.AmountOutAddr.String())
		frontRunGasCost := new(big.Int).SetUint64(tf.FrontRunTrade.TotalGasCost)
		sandwichRunGasCost := new(big.Int).SetUint64(tf.SandwichTrade.TotalGasCost)
		totalSandwichTradeGasCost := new(big.Int).Add(frontRunGasCost, sandwichRunGasCost)
		fmt.Println("total gas cost", totalSandwichTradeGasCost.String())

		endingTokenBalance, err := u.Web3Client.ReadERC20TokenBalance(ctx, tf.SandwichTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
		if err != nil {
			return err
		}
		profitTokenBalance := new(big.Int).Sub(endingTokenBalance, tf.FrontRunTrade.AmountIn)
		fmt.Println("profitTokenBalance", profitTokenBalance.String())
		fmt.Println("sandwichCalculatedProfit", tf.SandwichPrediction.ExpectedProfit.String())
		if profitTokenBalance.String() != tf.SandwichPrediction.ExpectedProfit.String() {
			return errors.New("profit token balance mismatch")
		}
	}
	return nil
}

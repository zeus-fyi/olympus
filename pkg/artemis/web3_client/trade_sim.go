package web3_client

import (
	"fmt"
	"math/big"
)

// TODO verify this is correct
func (u *UniswapV2Client) TradeSim(tf TradeExecutionFlow) (*big.Int, *big.Int) {
	sellAmount := big.NewInt(0)
	maxProfit := big.NewInt(0)
	if tf.TradeMethod == "swapExactTokensForETH" {
		startOffset := big.NewInt(1)
		offset := startOffset
		for {
			pair := tf.InitialPair
			diff := new(big.Int).Sub(offset, tf.UserTrade.AmountIn)
			if diff.Cmp(big.NewInt(0)) == 1 {
				break
			}
			amountOut := pair.PriceImpact(tf.UserTrade.AmountInAddr, offset)
			pair.PriceImpact(tf.UserTrade.AmountInAddr, tf.UserTrade.AmountIn)
			revenue := pair.PriceImpact(tf.SandwichTrade.AmountInAddr, amountOut.AmountOut)
			profit := new(big.Int).Sub(revenue.AmountOut, offset)
			if profit.Cmp(maxProfit) == 1 {
				maxProfit = profit
				sellAmount = offset
			}
			offset = offset.Add(offset, startOffset)
		}
	}
	return sellAmount, maxProfit
}

func (u *UniswapV2Client) TradeSimStep(tf TradeExecutionFlow) *big.Int {
	pair := tf.InitialPair
	//fmt.Println("frontRunAmountIn", tf.FrontRunTrade.AmountIn)
	frontRunTradeOutcome := pair.PriceImpact(tf.FrontRunTrade.AmountInAddr, tf.FrontRunTrade.AmountIn)
	pair.PriceImpact(tf.UserTrade.AmountInAddr, tf.UserTrade.AmountIn)
	//fmt.Println("userTradeAmountIn", tf.UserTrade.AmountIn.String())
	//fmt.Println("frontRunTradeOutcome", frontRunTradeOutcome.AmountOutAddr.String())
	//fmt.Println("frontRunTradeOutcome", tf.SandwichTrade.AmountInAddr.String())
	revenue := pair.PriceImpact(tf.SandwichTrade.AmountInAddr, tf.SandwichTrade.AmountIn)
	fmt.Println("revenue", revenue.AmountOut.String())
	maxProfit := new(big.Int).Sub(revenue.AmountOut, frontRunTradeOutcome.AmountIn)
	return maxProfit
}

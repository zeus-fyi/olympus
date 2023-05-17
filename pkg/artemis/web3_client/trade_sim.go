package web3_client

import (
	"fmt"
	"math/big"
)

func (u *UniswapV2Client) TradeSim(tf TradeExecutionFlow) *big.Int {
	pair := tf.InitialPair
	maxProfit := big.NewInt(0)
	if tf.TradeMethod == "swapExactTokensForETH" {
		startOffset := tf.FrontRunTrade.AmountIn.Div(tf.FrontRunTrade.AmountIn, big.NewInt(10000))
		offset := startOffset
		for {
			if offset.Cmp(tf.FrontRunTrade.AmountIn) == 1 {
				break
			}
			amountIn := new(big.Int).Add(tf.UserTrade.AmountIn, offset)
			amountOut := pair.PriceImpact(tf.UserTrade.AmountInAddr, tf.FrontRunTrade.AmountIn)
			pair.PriceImpact(tf.UserTrade.AmountInAddr, tf.FrontRunTrade.AmountIn)
			revenue := pair.PriceImpact(tf.UserTrade.AmountOutAddr, amountOut.AmountOut)
			profit := new(big.Int).Sub(revenue.AmountOut, amountIn)
			if profit.Cmp(maxProfit) == 1 {
				maxProfit = profit
			}
			offset = offset.Add(offset, startOffset)
		}
	}
	return maxProfit
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

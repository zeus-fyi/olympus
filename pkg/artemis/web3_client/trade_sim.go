package web3_client

import "math/big"

//func (u *UniswapV2Client) TradeSim(tf TradeExecutionFlow) (*big.Int, *big.Int) {
//	sellAmount := big.NewInt(0)
//	maxProfit := big.NewInt(0)
//	startOffset := big.NewInt(1000000000000)
//	offset := startOffset
//	for {
//		pair := tf.InitialPair
//		if offset.Cmp(tf.UserTrade.AmountIn) == 1 {
//			break
//		}
//		amountOut, _ := pair.PriceImpact(tf.UserTrade.AmountInAddr, offset)
//		pair.PriceImpact(tf.UserTrade.AmountInAddr, tf.UserTrade.AmountIn)
//		revenue, _ := pair.PriceImpact(tf.SandwichTrade.AmountInAddr, amountOut.AmountOut)
//		profit := new(big.Int).Sub(revenue.AmountOut, amountOut.AmountIn)
//		if profit.Cmp(maxProfit) == 1 && profit.Cmp(big.NewInt(0)) == 1 {
//			maxProfit = profit
//			sellAmount = offset
//		}
//		offset = new(big.Int).Add(offset, startOffset)
//	}
//	return sellAmount, maxProfit
//}

func (u *UniswapV2Client) TradeSimStep(tf TradeExecutionFlowInBigInt) *big.Int {
	pair := tf.InitialPair
	//fmt.Println("frontRunAmountIn", tf.FrontRunTrade.AmountIn)
	frontRunTradeOutcome, _ := pair.PriceImpact(tf.FrontRunTrade.AmountInAddr, tf.FrontRunTrade.AmountIn)
	pair.PriceImpact(tf.UserTrade.AmountInAddr, tf.UserTrade.AmountIn)
	//fmt.Println("userTradeAmountIn", tf.UserTrade.AmountIn.String())
	//fmt.Println("frontRunTradeOutcome", frontRunTradeOutcome.AmountOutAddr.String())
	//fmt.Println("frontRunTradeOutcome", tf.SandwichTrade.AmountInAddr.String())
	revenue, _ := pair.PriceImpact(tf.SandwichTrade.AmountInAddr, tf.SandwichTrade.AmountIn)
	//fmt.Println("revenue", revenue.AmountOut.String())
	maxProfit := new(big.Int).Sub(revenue.AmountOut, frontRunTradeOutcome.AmountIn)
	return maxProfit
}

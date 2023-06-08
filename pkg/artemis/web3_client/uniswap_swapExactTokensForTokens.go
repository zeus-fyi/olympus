package web3_client

import (
	"fmt"
)

func (u *UniswapClient) SwapExactTokensForTokens(tx MevTx, args map[string]interface{}) {
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		return
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		return
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		return
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		return
	}
	st := SwapExactTokensForTokensParams{
		AmountIn:     amountIn,
		AmountOutMin: amountOutMin,
		Path:         path,
		To:           to,
		Deadline:     deadline,
	}
	pd, err := u.GetPricingData(ctx, path)
	if err != nil {
		return
	}
	initialPair := pd.v2Pair
	tf := st.BinarySearch(pd.v2Pair)
	tf.InitialPair = initialPair.ConvertToJSONType()
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapExactTokensForTokens==================================")
		ts := &TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tf,
			TokenAddr:     path[0].String(),
			BuyWithAmount: st.AmountIn,
			MinimumAmount: st.AmountOutMin,
		}
		u.PrintTradeSummaries(ts)
		// u.PrintTradeSummaries(tx, tf, pd.v2Pair, path[0].String(), st.AmountIn, st.AmountOutMin)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell BuyWithAmount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapExactTokensForTokens==================================")
	}
	u.SwapExactTokensForTokensParamsSlice = append(u.SwapExactTokensForTokensParamsSlice, st)
}

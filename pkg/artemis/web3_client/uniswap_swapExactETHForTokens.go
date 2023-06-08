package web3_client

import (
	"fmt"
	"math/big"
)

func (u *UniswapClient) SwapExactETHForTokens(tx MevTx, args map[string]interface{}, payableEth *big.Int) {
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
	st := SwapExactETHForTokensParams{
		AmountOutMin: amountOutMin,
		Path:         path,
		To:           to,
		Deadline:     deadline,
		Value:        payableEth,
	}
	pd, err := u.GetPricingData(ctx, path)
	if err != nil {
		return
	}
	initialPair := pd.v2Pair
	tf := st.BinarySearch(pd.v2Pair)
	tf.InitialPair = initialPair.ConvertToJSONType()
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapExactETHForTokens==================================")
		ts := &TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tf,
			TokenAddr:     path[0].String(),
			BuyWithAmount: st.Value,
			MinimumAmount: st.AmountOutMin,
		}
		u.PrintTradeSummaries(ts)
		//u.PrintTradeSummaries(tx, tf, pd.v2Pair, path[0].String(), st.Value, st.AmountOutMin)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell BuyWithAmount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapExactETHForTokens==================================")
	}
	u.SwapExactETHForTokensParamsSlice = append(u.SwapExactETHForTokensParamsSlice, st)
}

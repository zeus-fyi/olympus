package web3_client

import (
	"fmt"
)

func (u *UniswapClient) SwapTokensForExactTokens(tx MevTx, args map[string]interface{}) {
	amountOut, err := ParseBigInt(args["amountOut"])
	if err != nil {
		return
	}
	amountInMax, err := ParseBigInt(args["amountInMax"])
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
	st := SwapTokensForExactTokensParams{
		AmountOut:   amountOut,
		AmountInMax: amountInMax,
		Path:        path,
		To:          to,
		Deadline:    deadline,
	}
	pd, err := u.GetPricingData(ctx, path)
	if err != nil {
		return
	}
	initialPair := pd.v2Pair
	tf := st.BinarySearch(pd.v2Pair)
	tf.InitialPair = initialPair.ConvertToJSONType()
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapTokensForExactTokens==================================")
		ts := TradeSummary{
			Tx:        tx,
			Tf:        tf,
			Pd:        pd,
			TokenAddr: path[0].String(),
			Amount:    st.AmountInMax,
			AmountMin: st.AmountOut,
		}
		u.PrintTradeSummaries2(&ts)
		//u.PrintTradeSummaries(tx, tf, pd.v2Pair, path[0].String(), st.AmountInMax, st.AmountOut)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapTokensForExactTokens==================================")
	}
	u.SwapTokensForExactTokensParamsSlice = append(u.SwapTokensForExactTokensParamsSlice, st)
}

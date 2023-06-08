package web3_client

import (
	"fmt"
	"math/big"
)

func (u *UniswapClient) SwapETHForExactTokens(tx MevTx, args map[string]interface{}, payableEth *big.Int) {
	amountOut, err := ParseBigInt(args["amountOut"])
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
	st := SwapETHForExactTokensParams{
		AmountOut: amountOut,
		Path:      path,
		To:        to,
		Deadline:  deadline,
		Value:     payableEth,
	}
	pd, err := u.GetPricingData(ctx, path)
	if err != nil {
		return
	}
	initialPair := pd.v2Pair
	tf := st.BinarySearch(pd.v2Pair)
	tf.InitialPair = initialPair.ConvertToJSONType()
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapETHForExactTokens==================================")
		ts := TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tf,
			TokenAddr:     path[0].String(),
			BuyWithAmount: st.Value,
			MinimumAmount: st.AmountOut,
		}
		u.PrintTradeSummaries(&ts)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell BuyWithAmount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapETHForExactTokens==================================")
	}
	u.SwapETHForExactTokensParamsSlice = append(u.SwapETHForExactTokensParamsSlice, st)
}

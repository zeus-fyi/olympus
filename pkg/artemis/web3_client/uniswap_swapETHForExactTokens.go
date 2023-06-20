package web3_client

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

type SwapETHForExactTokensParams struct {
	AmountOut *big.Int           `json:"amountOut"`
	Path      []accounts.Address `json:"path"`
	To        accounts.Address   `json:"to"`
	Deadline  *big.Int           `json:"deadline"`
	Value     *big.Int           `json:"value"`
}

type JSONSwapETHForExactTokensParams struct {
	AmountOut string             `json:"amountOut"`
	Path      []accounts.Address `json:"path"`
	To        accounts.Address   `json:"to"`
	Deadline  string             `json:"deadline"`
	Value     string             `json:"value"`
}

func (s *SwapETHForExactTokensParams) ConvertToJSONType() *JSONSwapETHForExactTokensParams {
	return &JSONSwapETHForExactTokensParams{
		AmountOut: s.AmountOut.String(),
		Path:      s.Path,
		To:        s.To,
		Deadline:  s.Deadline.String(),
		Value:     s.Value.String(),
	}
}
func (s *SwapETHForExactTokensParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlowJSON {
	// Value == variable
	// AmountOut == required for trade
	low := big.NewInt(0)
	high := new(big.Int).Set(s.Value)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlowJSON{
		Trade: Trade{
			TradeMethod:                     "swapETHForExactTokens",
			JSONSwapETHForExactTokensParams: s.ConvertToJSONType(),
		},
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		mid = new(big.Int).Add(low, high)
		mid = mid.Div(mid, big.NewInt(2))
		// Front run trade
		toFrontRun, err := mockPairResp.PriceImpact(s.Path[0], mid)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		// User trade
		to, err := mockPairResp.PriceImpact(s.Path[0], s.Value)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOut)
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(mid, big.NewInt(1))
			continue
		}
		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich, err := mockPairResp.PriceImpact(s.Path[1], sandwichDump)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = toFrontRun.ConvertToJSONType()
			tf.UserTrade = to.ConvertToJSONType()
			tf.SandwichTrade = toSandwich.ConvertToJSONType()
		}
		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(mid, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(mid, big.NewInt(1))
		}
	}
	sp := SandwichTradePrediction{
		SellAmount:     tokenSellAmountAtMaxProfit,
		ExpectedProfit: maxProfit,
	}
	tf.SandwichPrediction = sp.ConvertToJSONType()
	return tf
}

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
		fmt.Println("txHash: ", tx.Tx.Hash().String())
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell BuyWithAmount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapETHForExactTokens==================================")
	}
	u.SwapETHForExactTokensParamsSlice = append(u.SwapETHForExactTokensParamsSlice, st)
}

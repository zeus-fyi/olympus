package web3_client

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
)

type SwapTokensForExactTokensParams struct {
	AmountOut   *big.Int           `json:"amountOut"`
	AmountInMax *big.Int           `json:"amountInMax"`
	Path        []accounts.Address `json:"path"`
	To          accounts.Address   `json:"to"`
	Deadline    *big.Int           `json:"deadline"`
}

type JSONSwapTokensForExactTokensParams struct {
	AmountOut   string             `json:"amountOut"`
	AmountInMax string             `json:"amountInMax"`
	Path        []accounts.Address `json:"path"`
	To          accounts.Address   `json:"to"`
	Deadline    string             `json:"deadline"`
}

func (s *SwapTokensForExactTokensParams) ConvertToJSONType() *JSONSwapTokensForExactTokensParams {
	return &JSONSwapTokensForExactTokensParams{
		AmountOut:   s.AmountOut.String(),
		AmountInMax: s.AmountInMax.String(),
		Path:        s.Path,
		To:          s.To,
		Deadline:    s.Deadline.String(),
	}
}

func (s *SwapTokensForExactTokensParams) BinarySearch(pair uniswap_pricing.UniswapV2Pair) (TradeExecutionFlow, error) {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountInMax)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		InitialPair: &pair,
		Trade: Trade{
			TradeMethod:                        swapTokensForExactTokens,
			JSONSwapTokensForExactTokensParams: s.ConvertToJSONType(),
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
			return tf, err
		}
		// User trade
		to, err := mockPairResp.PriceImpact(s.Path[0], s.AmountInMax)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf, err
		}
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOut)
		// if diff <= 0 then it searches left
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(mid, big.NewInt(1))
			continue
		}
		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich, err := mockPairResp.PriceImpact(s.Path[1], sandwichDump)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf, err
		}
		profit := artemis_eth_units.SubBigInt(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = toFrontRun
			tf.UserTrade = to
			tf.SandwichTrade = toSandwich
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
	tf.SandwichPrediction = sp
	return tf, nil
}

func (s *SwapTokensForExactTokensParams) Decode(args map[string]interface{}) error {
	amountOut, err := ParseBigInt(args["amountOut"])
	if err != nil {
		log.Warn().Msg("SwapTokensForExactTokensParams: error parsing amountOut")
		return err
	}
	amountInMax, err := ParseBigInt(args["amountInMax"])
	if err != nil {
		log.Warn().Msg("SwapTokensForExactTokensParams: error parsing amountInMax")
		return err
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		log.Warn().Msg("SwapTokensForExactTokensParams: error parsing path")
		return err
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		log.Warn().Msg("SwapTokensForExactTokensParams: error parsing to")
		return err
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		log.Warn().Msg("SwapTokensForExactTokensParams: error parsing deadline")
		return err
	}
	s.AmountOut = amountOut
	s.AmountInMax = amountInMax
	s.Path = path
	s.To = to
	s.Deadline = deadline
	return nil
}

func (u *UniswapClient) SwapTokensForExactTokens(tx MevTx, args map[string]interface{}) error {
	st := SwapTokensForExactTokensParams{}
	err := st.Decode(args)
	if err != nil {
		return err
	}
	path := st.Path
	pd, err := u.GetV2PricingData(ctx, path)
	if err != nil {
		return err
	}
	tf, err := st.BinarySearch(pd.V2Pair)
	if err != nil {
		return err
	}
	tfJSON, err := tf.ConvertToJSONType()
	if err != nil {
		return err
	}
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapTokensForExactTokens==================================")
		ts := TradeSummary{
			Tx:            tx,
			Tf:            tfJSON,
			Pd:            pd,
			TokenAddr:     path[0].String(),
			BuyWithAmount: st.AmountInMax,
			MinimumAmount: st.AmountOut,
		}
		u.PrintTradeSummaries(&ts)
		//u.PrintTradeSummaries(tx, tf, pd.V2Pair, path[0].String(), st.AmountInMax, st.AmountOut)
		fmt.Println("txHash: ", tx.Tx.Hash().String())
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "SandwichPrediction Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapTokensForExactTokens==================================")
	}
	u.SwapTokensForExactTokensParamsSlice = append(u.SwapTokensForExactTokensParamsSlice, st)
	return nil
}

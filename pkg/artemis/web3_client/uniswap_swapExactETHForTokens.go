package web3_client

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
)

type SwapExactETHForTokensParams struct {
	AmountOutMin *big.Int           `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Value        *big.Int           `json:"value"`
	Deadline     *big.Int           `json:"deadline"`
}

type JSONSwapExactETHForTokensParams struct {
	AmountOutMin string             `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Value        string             `json:"value"`
	Deadline     string             `json:"deadline"`
}

func (s *SwapExactETHForTokensParams) ConvertToJSONType() *JSONSwapExactETHForTokensParams {
	return &JSONSwapExactETHForTokensParams{
		AmountOutMin: s.AmountOutMin.String(),
		Path:         s.Path,
		To:           s.To,
		Value:        s.Value.String(),
		Deadline:     s.Deadline.String(),
	}
}

func (s *JSONSwapExactETHForTokensParams) ConvertToBigIntType() *SwapExactETHForTokensParams {
	return &SwapExactETHForTokensParams{
		AmountOutMin: artemis_eth_units.NewBigIntFromStr(s.AmountOutMin),
		Path:         s.Path,
		To:           s.To,
		Value:        artemis_eth_units.NewBigIntFromStr(s.Value),
		Deadline:     artemis_eth_units.NewBigIntFromStr(s.Deadline),
	}
}

func (s *SwapExactETHForTokensParams) BinarySearch(pair uniswap_pricing.UniswapV2Pair) (TradeExecutionFlow, error) {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.Value)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		InitialPair: &pair,
		Trade: Trade{
			TradeMethod:                     swapExactETHForTokens,
			JSONSwapExactETHForTokensParams: s.ConvertToJSONType(),
		},
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		mid = new(big.Int).Add(low, high)
		mid = DivideByHalf(mid)
		// Front run trade
		toFrontRun, err := mockPairResp.PriceImpact(s.Path[0], mid)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf, err
		}
		// User trade
		to, err := mockPairResp.PriceImpact(s.Path[0], s.Value)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf, err
		}
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOutMin)
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

func (s *SwapExactETHForTokensParams) Decode(args map[string]interface{}, payableEth *big.Int) error {
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		log.Warn().Msg("SwapExactETHForTokensParams: error in parsing amountOutMin")
		return err
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		log.Warn().Msg("SwapExactETHForTokensParams: error in parsing path")
		return err
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		log.Warn().Msg("SwapExactETHForTokensParams: error in parsing to")
		return err
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		log.Warn().Msg("SwapExactETHForTokensParams: error in parsing deadline")
		return err
	}
	s.AmountOutMin = amountOutMin
	s.Path = path
	s.To = to
	s.Deadline = deadline
	s.Value = payableEth
	return nil
}

func (u *UniswapClient) SwapExactETHForTokens(tx MevTx, args map[string]interface{}, payableEth *big.Int) error {
	st := SwapExactETHForTokensParams{}
	err := st.Decode(args, payableEth)
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
		fmt.Println("\nsandwich: ==================================SwapExactETHForTokens==================================")
		ts := &TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tfJSON,
			TokenAddr:     path[0].String(),
			BuyWithAmount: st.Value,
			MinimumAmount: st.AmountOutMin,
		}
		u.PrintTradeSummaries(ts)
		//u.PrintTradeSummaries(tx, tf, pd.V2Pair, path[0].String(), st.Value, st.AmountOutMin)
		fmt.Println("txHash: ", tx.Tx.Hash().String())
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "SandwichPrediction Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapExactETHForTokens==================================")
	}
	u.SwapExactETHForTokensParamsSlice = append(u.SwapExactETHForTokensParamsSlice, st)
	return nil
}

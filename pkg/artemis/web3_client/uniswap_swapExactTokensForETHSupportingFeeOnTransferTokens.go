package web3_client

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
)

type SwapExactTokensForETHSupportingFeeOnTransferTokensParams struct {
	AmountIn     *big.Int           `json:"amountIn"`
	AmountOutMin *big.Int           `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     *big.Int           `json:"deadline"`
}

type JSONSwapExactTokensForETHSupportingFeeOnTransferTokensParams struct {
	AmountIn     string             `json:"amountIn"`
	AmountOutMin string             `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     string             `json:"deadline"`
}

func (s *JSONSwapExactTokensForETHSupportingFeeOnTransferTokensParams) ConvertToBigIntType() *SwapExactTokensForETHSupportingFeeOnTransferTokensParams {
	return &SwapExactTokensForETHSupportingFeeOnTransferTokensParams{
		AmountIn:     artemis_eth_units.NewBigIntFromStr(s.AmountIn),
		AmountOutMin: artemis_eth_units.NewBigIntFromStr(s.AmountOutMin),
		Path:         s.Path,
		To:           s.To,
		Deadline:     artemis_eth_units.NewBigIntFromStr(s.Deadline),
	}
}

func (s *SwapExactTokensForETHSupportingFeeOnTransferTokensParams) ConvertToJSONType() *JSONSwapExactTokensForETHSupportingFeeOnTransferTokensParams {
	return &JSONSwapExactTokensForETHSupportingFeeOnTransferTokensParams{
		AmountIn:     s.AmountIn.String(),
		AmountOutMin: s.AmountOutMin.String(),
		Path:         s.Path,
		To:           s.To,
		Deadline:     s.Deadline.String(),
	}
}

func (u *UniswapClient) SwapExactTokensForETHSupportingFeeOnTransferTokens(tx MevTx, args map[string]interface{}) error {
	st := SwapExactTokensForETHSupportingFeeOnTransferTokensParams{}
	err := st.Decode(args)
	if err != nil {
		log.Error().Err(err).Msg("error in decoding")
		return err
	}
	path := st.Path
	pd, err := u.GetV2PricingData(ctx, path)
	if err != nil {
		return err
	}
	initialPair := pd.V2Pair
	tf, err := st.BinarySearch(pd.V2Pair)
	if err != nil {
		log.Error().Err(err).Msg("error in binary search")
		return err
	}
	tf.InitialPair = initialPair.ConvertToJSONType()
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapExactTokensForETHSupportingFeeOnTransferTokens==================================")
		ts := &TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tf,
			TokenAddr:     path[0].String(),
			BuyWithAmount: st.AmountIn,
			MinimumAmount: st.AmountOutMin,
		}
		u.PrintTradeSummaries(ts)
		//u.PrintTradeSummaries(tx, tf, pd.V2Pair, path[0].String(), st.Value, st.AmountOutMin)
		fmt.Println("txHash: ", tx.Tx.Hash().String())
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "SandwichPrediction Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapExactTokensForETHSupportingFeeOnTransferTokens==================================")
	}
	u.SwapExactTokensForETHSupportingFeeOnTransferTokensParamsSlice = append(u.SwapExactTokensForETHSupportingFeeOnTransferTokensParamsSlice, st)
	return nil
}

func (s *SwapExactTokensForETHSupportingFeeOnTransferTokensParams) Decode(args map[string]interface{}) error {
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForETHSupportingFeeOnTransferTokensParams: error in parsing amountIn")
		return err
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForETHSupportingFeeOnTransferTokensParams: error in parsing amountOutMin")
		return err
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForETHSupportingFeeOnTransferTokensParams: error in parsing path")
		return err
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForETHSupportingFeeOnTransferTokensParams: error in parsing to")
		return err
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForETHSupportingFeeOnTransferTokensParams: error in parsing deadline")
		return err
	}
	s.AmountIn = amountIn
	s.AmountOutMin = amountOutMin
	s.Path = path
	s.To = to
	s.Deadline = deadline
	return nil
}

const (
	swapExactTokensForETHSupportingFeeOnTransferTokensMoniker = "swapExactTokensForETHSupportingFeeOnTransferTokens"
)

func (s *SwapExactTokensForETHSupportingFeeOnTransferTokensParams) BinarySearch(pair uniswap_pricing.UniswapV2Pair) (TradeExecutionFlowJSON, error) {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlowJSON{
		Trade: Trade{
			TradeMethod: swapExactTokensForETHSupportingFeeOnTransferTokensMoniker,
			JSONSwapExactTokensForETHSupportingFeeOnTransferTokensParams: s.ConvertToJSONType(),
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
		to, err := mockPairResp.PriceImpact(s.Path[0], s.AmountIn)
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
	return tf, nil
}

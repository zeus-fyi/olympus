package web3_client

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
)

const (
	swapExactTokensForTokensSupportingFeeOnTransferTokens = "swapExactTokensForTokensSupportingFeeOnTransferTokens"
)

type SwapExactTokensForTokensSupportingFeeOnTransferTokensParams struct {
	AmountIn     *big.Int           `json:"amountIn"`
	AmountOutMin *big.Int           `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     *big.Int           `json:"deadline"`
}

func (s *SwapExactTokensForTokensSupportingFeeOnTransferTokensParams) BinarySearch(pair uniswap_pricing.UniswapV2Pair) (TradeExecutionFlowJSON, error) {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlowJSON{
		Trade: Trade{
			TradeMethod: swapExactTokensForTokensSupportingFeeOnTransferTokens,
			JSONSwapExactTokensForTokensSupportingFeeOnTransferTokensParams: s.ConvertToJSONType(),
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
func (s *SwapExactTokensForTokensSupportingFeeOnTransferTokensParams) Decode(args map[string]interface{}) error {
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForTokensSupportingFeeOnTransferTokensParams: error parsing amountIn")
		return err
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForTokensSupportingFeeOnTransferTokensParams: error parsing amountOutMin")
		return err
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForTokensSupportingFeeOnTransferTokensParams: error parsing path")
		return err
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForTokensSupportingFeeOnTransferTokensParams: error parsing to")
		return err
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForTokensSupportingFeeOnTransferTokensParams: error parsing deadline")
		return err
	}
	s.AmountIn = amountIn
	s.AmountOutMin = amountOutMin
	s.Path = path
	s.To = to
	s.Deadline = deadline
	return nil
}

func (u *UniswapClient) SwapExactTokensForTokensSupportingFeeOnTransferTokens(tx MevTx, args map[string]interface{}) error {
	st := SwapExactTokensForTokensSupportingFeeOnTransferTokensParams{}
	err := st.Decode(args)
	if err != nil {
		return err
	}
	pd, err := u.GetV2PricingData(ctx, st.Path)
	if err != nil {
		return err
	}
	path := st.Path
	initialPair := pd.V2Pair
	tf, err := st.BinarySearch(pd.V2Pair)
	if err != nil {
		return err
	}
	tf.InitialPair = initialPair.ConvertToJSONType()
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapExactTokensForTokensSupportingFeeOnTransferTokens==================================")
		ts := &TradeSummary{
			Tx:            tx,
			Pd:            pd,
			Tf:            tf,
			TokenAddr:     path[0].String(),
			BuyWithAmount: st.AmountIn,
			MinimumAmount: st.AmountOutMin,
		}
		u.PrintTradeSummaries(ts)
		fmt.Println("txHash: ", tx.Tx.Hash().String())
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "SandwichPrediction Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapExactTokensForTokensSupportingFeeOnTransferTokens==================================")
	}
	u.SwapExactTokensForTokensSupportingFeeOnTransferTokensSlice = append(u.SwapExactTokensForTokensSupportingFeeOnTransferTokensSlice, st)
	return nil
}

type JSONSwapExactTokensForTokensSupportingFeeOnTransferTokensParams struct {
	AmountIn     string             `json:"amountIn"`
	AmountOutMin string             `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     string             `json:"deadline"`
}

func (s *SwapExactTokensForTokensSupportingFeeOnTransferTokensParams) ConvertToJSONType() *JSONSwapExactTokensForTokensSupportingFeeOnTransferTokensParams {
	return &JSONSwapExactTokensForTokensSupportingFeeOnTransferTokensParams{
		AmountIn:     s.AmountIn.String(),
		AmountOutMin: s.AmountOutMin.String(),
		Path:         s.Path,
		To:           s.To,
		Deadline:     s.Deadline.String(),
	}
}

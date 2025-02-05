package web3_client

import (
	"context"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
)

type SwapExactTokensForTokensParamsV3 struct {
	AmountIn     *big.Int           `json:"amountIn"`
	AmountOutMin *big.Int           `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
}

type JSONSwapExactTokensForTokensParamsV3 struct {
	AmountIn     string             `json:"amountIn"`
	AmountOutMin string             `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
}

func (s *SwapExactTokensForTokensParamsV3) BinarySearch(pair uniswap_pricing.UniswapV2Pair) (TradeExecutionFlow, error) {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		InitialPair: &pair,
		Trade: Trade{
			TradeMethod:                          swapExactTokensForTokens,
			JSONSwapExactTokensForTokensParamsV3: s.ConvertToJSONType(),
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

func (s *SwapExactTokensForTokensParamsV3) Decode(ctx context.Context, args map[string]interface{}) error {
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForTokensParamsV3: error parsing amountIn")
		return err
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForTokensParamsV3: error parsing amountOutMin")
		return err
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForTokensParamsV3: error parsing path")
		return err
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		log.Warn().Msg("SwapExactTokensForTokensParamsV3: error parsing to")
		return err
	}
	s.AmountIn = amountIn
	s.AmountOutMin = amountOutMin
	s.Path = path
	s.To = to
	return nil
}

func (s *SwapExactTokensForTokensParamsV3) ConvertToJSONType() *JSONSwapExactTokensForTokensParamsV3 {
	return &JSONSwapExactTokensForTokensParamsV3{
		AmountIn:     s.AmountIn.String(),
		AmountOutMin: s.AmountOutMin.String(),
		Path:         s.Path,
		To:           s.To,
	}
}

package web3_client

import (
	"math/big"

	"github.com/rs/zerolog/log"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

func (s *V3SwapExactOutParams) BinarySearch(pd *PricingData) TradeExecutionFlowJSON {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountInMax)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlowJSON{
		Trade: Trade{
			TradeMethod:              V3SwapExactOut,
			JSONV3SwapExactOutParams: s.ConvertToJSONType(),
		},
	}
	frontRunTokenIn := pd.V3Pair.Token0
	sandwichTokenIn := pd.V3Pair.Token1
	if s.Path.TokenIn.Hex() == pd.V3Pair.Token1.Address.Hex() {
		frontRunTokenIn = pd.V3Pair.Token1
		sandwichTokenIn = pd.V3Pair.Token0
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pd.V3Pair
		mid = new(big.Int).Add(low, high)
		mid = DivideByHalf(mid)
		// Front run trade
		amountInFrontRun := mid
		toFrontRun, _, err := mockPairResp.PriceImpact(ctx, frontRunTokenIn, amountInFrontRun)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		// User trade
		userTrade, _, err := mockPairResp.PriceImpact(ctx, frontRunTokenIn, s.AmountInMax)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		difference := new(big.Int).Sub(userTrade.Quotient(), s.AmountOut)
		// if diff <= 0 then it searches left
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(mid, big.NewInt(1))
			continue
		}
		// Sandwich trade
		toSandwich, _, err := mockPairResp.PriceImpact(ctx, sandwichTokenIn, toFrontRun.Quotient())
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		profit := new(big.Int).Sub(toSandwich.Quotient(), toFrontRun.Quotient())
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = artemis_trading_types.JSONTradeOutcome{
				AmountIn:      amountInFrontRun.String(),
				AmountInAddr:  frontRunTokenIn.Address,
				AmountOut:     toFrontRun.Quotient().String(),
				AmountOutAddr: sandwichTokenIn.Address,
			}
			tf.UserTrade = artemis_trading_types.JSONTradeOutcome{
				AmountIn:      s.AmountInMax.String(),
				AmountInAddr:  frontRunTokenIn.Address,
				AmountOut:     userTrade.Quotient().String(),
				AmountOutAddr: sandwichTokenIn.Address,
			}
			tf.SandwichTrade = artemis_trading_types.JSONTradeOutcome{
				AmountIn:      toFrontRun.Quotient().String(),
				AmountInAddr:  sandwichTokenIn.Address,
				AmountOut:     toSandwich.Quotient().String(),
				AmountOutAddr: frontRunTokenIn.Address,
			}
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

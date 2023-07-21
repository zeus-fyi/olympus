package web3_client

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

type SwapExactInputSingleArgs struct {
	TokenIn           accounts.Address `json:"tokenIn"`
	AmountIn          *big.Int         `json:"amountIn"`
	TokenOut          accounts.Address `json:"tokenOut"`
	AmountOutMinimum  *big.Int         `json:"amountOutMinimum"`
	Fee               *big.Int         `json:"fee"`
	Recipient         accounts.Address `json:"recipient"`
	SqrtPriceLimitX96 *big.Int         `json:"sqrtPriceLimitX96"`

	TokenFeePath artemis_trading_types.TokenFeePath `json:"tokenFeePath,omitempty"`
}

type JSONSwapExactInputSingleArgs struct {
	TokenIn           accounts.Address                   `json:"tokenIn"`
	AmountIn          string                             `json:"amountIn"`
	TokenOut          accounts.Address                   `json:"tokenOut"`
	AmountOutMinimum  string                             `json:"amountOutMinimum"`
	Fee               string                             `json:"fee"`
	Recipient         accounts.Address                   `json:"recipient"`
	SqrtPriceLimitX96 string                             `json:"sqrtPriceLimitX96"`
	TokenFeePath      artemis_trading_types.TokenFeePath `json:"tokenFeePath,omitempty"`
}

func (s *SwapExactInputSingleArgs) BinarySearch(pd *uniswap_pricing.UniswapPricingData) (TradeExecutionFlow, error) {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		InitialPairV3: &pd.V3Pair,
		Trade: Trade{
			TradeMethod:                  swapExactInputSingle,
			JSONSwapExactInputSingleArgs: s.ConvertToJSONType(),
		},
	}
	frontRunTokenIn := pd.V3Pair.Token0
	sandwichTokenIn := pd.V3Pair.Token1
	if s.TokenFeePath.TokenIn.Hex() == pd.V3Pair.Token1.Address.Hex() {
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
			return tf, err
		}
		// User trade
		userTrade, _, err := mockPairResp.PriceImpact(ctx, frontRunTokenIn, s.AmountIn)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf, err
		}
		difference := new(big.Int).Sub(userTrade.Quotient(), s.AmountOutMinimum)
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(mid, big.NewInt(1))
			continue
		}
		// Sandwich trade
		toSandwich, _, err := mockPairResp.PriceImpact(ctx, sandwichTokenIn, toFrontRun.Quotient())
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf, err
		}
		profit := new(big.Int).Sub(toSandwich.Quotient(), toFrontRun.Quotient())
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = artemis_trading_types.TradeOutcome{
				AmountIn:      amountInFrontRun,
				AmountInAddr:  frontRunTokenIn.Address,
				AmountOut:     toFrontRun.Quotient(),
				AmountOutAddr: sandwichTokenIn.Address,
			}
			tf.UserTrade = artemis_trading_types.TradeOutcome{
				AmountIn:      s.AmountIn,
				AmountInAddr:  frontRunTokenIn.Address,
				AmountOut:     userTrade.Quotient(),
				AmountOutAddr: sandwichTokenIn.Address,
			}
			tf.SandwichTrade = artemis_trading_types.TradeOutcome{
				AmountIn:      toFrontRun.Quotient(),
				AmountInAddr:  sandwichTokenIn.Address,
				AmountOut:     toSandwich.Quotient(),
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
	tf.SandwichPrediction = sp
	return tf, nil
}

func (s *JSONSwapExactInputSingleArgs) ConvertToBigIntType() (SwapExactInputSingleArgs, error) {
	amountIn, ok1 := new(big.Int).SetString(s.AmountIn, 10)
	amountOutMinimum, ok2 := new(big.Int).SetString(s.AmountOutMinimum, 10)
	fee, ok3 := new(big.Int).SetString(s.Fee, 10)
	sqrtPriceLimitX96, ok4 := new(big.Int).SetString(s.SqrtPriceLimitX96, 10)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		log.Error().Msg("error converting to big int type")
		return SwapExactInputSingleArgs{}, errors.New("error converting to big int type")
	}
	return SwapExactInputSingleArgs{
		TokenIn:           s.TokenIn,
		AmountIn:          amountIn,
		TokenOut:          s.TokenOut,
		AmountOutMinimum:  amountOutMinimum,
		Fee:               fee,
		Recipient:         s.Recipient,
		SqrtPriceLimitX96: sqrtPriceLimitX96,
		TokenFeePath:      s.TokenFeePath,
	}, nil
}

func (s *SwapExactInputSingleArgs) ConvertToJSONType() *JSONSwapExactInputSingleArgs {
	return &JSONSwapExactInputSingleArgs{
		TokenIn:           s.TokenIn,
		AmountIn:          s.AmountIn.String(),
		TokenOut:          s.TokenOut,
		AmountOutMinimum:  s.AmountOutMinimum.String(),
		Fee:               s.Fee.String(),
		Recipient:         s.Recipient,
		SqrtPriceLimitX96: s.SqrtPriceLimitX96.String(),
		TokenFeePath:      s.TokenFeePath,
	}
}

func (s *SwapExactInputSingleArgs) Decode(ctx context.Context, args map[string]interface{}) error {
	params, ok := args["params"].(struct {
		TokenIn           common.Address "json:\"tokenIn\""
		TokenOut          common.Address "json:\"tokenOut\""
		Fee               *big.Int       "json:\"fee\""
		Recipient         common.Address "json:\"recipient\""
		AmountIn          *big.Int       "json:\"amountIn\""
		AmountOutMinimum  *big.Int       "json:\"amountOutMinimum\""
		SqrtPriceLimitX96 *big.Int       "json:\"sqrtPriceLimitX96\""
	})
	if !ok {
		log.Warn().Msg("SwapExactInputSingleArgs: Decode params is not of the expected type")
		return fmt.Errorf("params is not of the expected type")
	}
	s.TokenIn = accounts.Address(params.TokenIn)
	s.TokenOut = accounts.Address(params.TokenOut)
	s.Fee = params.Fee
	s.Recipient = accounts.Address(params.Recipient)
	s.AmountIn = params.AmountIn
	s.AmountOutMinimum = params.AmountOutMinimum
	s.SqrtPriceLimitX96 = params.SqrtPriceLimitX96
	return nil
}

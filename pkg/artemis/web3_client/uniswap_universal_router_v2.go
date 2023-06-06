package web3_client

import (
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

/*
V2_SWAP_EXACT_IN
address The recipient of the output of the trade
uint256 The amount of input tokens for the trade
uint256 The minimum amount of output tokens the user wants
address[] The UniswapV2 token path to trade along
bool A flag for whether the input tokens should come from the msg.sender (through Permit2) or whether the funds are already in the UniversalRouter

V2_SWAP_EXACT_OUT
address The recipient of the output of the trade
uint256 The amount of output tokens to receive
uint256 The maximum number of input tokens that should be spent
address[] The UniswapV2 token path to trade along
bool A flag for whether the input tokens should come from the msg.sender (through Permit2) or whether the funds are already in the UniversalRouter
*/

const (
	V2SwapExactIn  = "V2_SWAP_EXACT_IN"
	V2SwapExactOut = "V2_SWAP_EXACT_OUT"
)

type V2SwapExactInParams struct {
	AmountIn        *big.Int           `json:"amountIn"`
	AmountOutMin    *big.Int           `json:"amountOutMin"`
	Path            []accounts.Address `json:"path"`
	To              accounts.Address   `json:"to"`
	InputFromSender bool               `json:"inputFromSender"`
}

type JSONV2SwapExactInParams struct {
	AmountIn        string             `json:"amountIn"`
	AmountOutMin    string             `json:"amountOutMin"`
	Path            []accounts.Address `json:"path"`
	To              accounts.Address   `json:"to"`
	InputFromSender bool               `json:"inputFromSender"`
}

func (s *JSONV2SwapExactInParams) ConvertToBigIntType() *V2SwapExactInParams {
	amountIn, _ := new(big.Int).SetString(s.AmountIn, 10)
	amountOutMin, _ := new(big.Int).SetString(s.AmountOutMin, 10)
	return &V2SwapExactInParams{
		AmountIn:        amountIn,
		AmountOutMin:    amountOutMin,
		Path:            s.Path,
		To:              s.To,
		InputFromSender: s.InputFromSender,
	}
}

func (s *V2SwapExactInParams) ConvertToJSONType() *JSONV2SwapExactInParams {
	return &JSONV2SwapExactInParams{
		AmountIn:        s.AmountIn.String(),
		AmountOutMin:    s.AmountOutMin.String(),
		Path:            s.Path,
		To:              s.To,
		InputFromSender: s.InputFromSender,
	}
}

type V2SwapExactOutParams struct {
	AmountInMax     *big.Int           `json:"amountInMax"`
	AmountOut       *big.Int           `json:"amountOut"`
	Path            []accounts.Address `json:"path"`
	To              accounts.Address   `json:"to"`
	InputFromSender bool               `json:"inputFromSender"`
}

type JSONV2SwapExactOutParams struct {
	AmountInMax     string             `json:"amountInMax"`
	AmountOut       string             `json:"amountOut"`
	Path            []accounts.Address `json:"path"`
	To              accounts.Address   `json:"to"`
	InputFromSender bool               `json:"inputFromSender"`
}

func (s *JSONV2SwapExactOutParams) ConvertToBigIntType() *V2SwapExactOutParams {
	amountInMax, _ := new(big.Int).SetString(s.AmountInMax, 10)
	amountOut, _ := new(big.Int).SetString(s.AmountOut, 10)
	return &V2SwapExactOutParams{
		AmountInMax:     amountInMax,
		AmountOut:       amountOut,
		Path:            s.Path,
		To:              s.To,
		InputFromSender: s.InputFromSender,
	}
}

func (s *V2SwapExactOutParams) ConvertToJSONType() *JSONV2SwapExactOutParams {
	return &JSONV2SwapExactOutParams{
		AmountInMax:     s.AmountInMax.String(),
		AmountOut:       s.AmountOut.String(),
		Path:            s.Path,
		To:              s.To,
		InputFromSender: s.InputFromSender,
	}
}

func (s *V2SwapExactOutParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountInMax)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		Trade: Trade{
			TradeMethod:              V2SwapExactOut,
			JSONV2SwapExactOutParams: s.ConvertToJSONType(),
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
			return tf
		}
		// User trade
		to, err := mockPairResp.PriceImpact(s.Path[0], s.AmountInMax)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
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

func (s *V2SwapExactInParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		Trade: Trade{
			TradeMethod:             V2SwapExactIn,
			JSONV2SwapExactInParams: s.ConvertToJSONType(),
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
			return tf
		}
		// User trade
		to, err := mockPairResp.PriceImpact(s.Path[0], s.AmountIn)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
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

package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

/*
The inputs for V3_SWAP_EXACT_IN is the encoding of 5 parameters:

address The recipient of the output of the trade
uint256 The amount of input tokens for the trade
uint256 The minimum amount of output tokens the user wants
bytes The UniswapV3 path you want to trade along
bool A flag for whether the input funds should come from the caller (through Permit2) or whether the funds are already in the UniversalRouter
*/

const (
	V3SwapExactIn  = "V3_SWAP_EXACT_IN"
	V3SwapExactOut = "V3_SWAP_EXACT_OUT"
)

type V3SwapExactInParams struct {
	AmountIn     *big.Int           `json:"amountIn"`
	AmountOutMin *big.Int           `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	PayerIsUser  bool               `json:"payerIsUser"`
}

type JSONV3SwapExactInParams struct {
	AmountIn     string             `json:"amountIn"`
	AmountOutMin string             `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	PayerIsUser  bool               `json:"payerIsUser"`
}

func (s *V3SwapExactInParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoder.Methods[V3SwapExactIn].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		fmt.Println(err)
	}
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		return err
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		return err
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return err
	}
	to, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	payerIsSender := args["payerIsUser"].(bool)
	s.AmountIn = amountIn
	s.AmountOutMin = amountOutMin
	s.Path = path
	s.To = to
	s.PayerIsUser = payerIsSender
	return err
}

func (u *UniswapClient) V3SwapExactIn(tx MevTx, args map[string]interface{}) {}

func (s *JSONV3SwapExactInParams) ConvertToBigIntType() *V3SwapExactInParams {
	amountIn, _ := new(big.Int).SetString(s.AmountIn, 10)
	amountOutMin, _ := new(big.Int).SetString(s.AmountOutMin, 10)
	return &V3SwapExactInParams{
		AmountIn:     amountIn,
		AmountOutMin: amountOutMin,
		Path:         s.Path,
		To:           s.To,
		PayerIsUser:  s.PayerIsUser,
	}
}

func (s *V3SwapExactInParams) ConvertToJSONType() *JSONV3SwapExactInParams {
	return &JSONV3SwapExactInParams{
		AmountIn:     s.AmountIn.String(),
		AmountOutMin: s.AmountOutMin.String(),
		Path:         s.Path,
		To:           s.To,
		PayerIsUser:  s.PayerIsUser,
	}
}

/*
V3_SWAP_EXACT_OUT
address The recipient of the output of the trade
uint256 The amount of output tokens to receive
uint256 The maximum number of input tokens that should be spent
bytes The UniswapV3 encoded path to trade along
bool A flag for whether the input tokens should come from the msg.sender (through Permit2) or whether the funds are already in the UniversalRouter
*/

type V3SwapExactOutParams struct {
	AmountInMax *big.Int           `json:"amountInMax"`
	AmountOut   *big.Int           `json:"amountOut"`
	Path        []accounts.Address `json:"path"`
	To          accounts.Address   `json:"to"`
	PayerIsUser bool               `json:"payerIsUser"`
}

type JSONV3SwapExactOutParams struct {
	AmountInMax string             `json:"amountInMax"`
	AmountOut   string             `json:"amountOut"`
	Path        []accounts.Address `json:"path"`
	To          accounts.Address   `json:"to"`
	PayerIsUser bool               `json:"payerIsUser"`
}

func (s *V3SwapExactOutParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoder.Methods[V3SwapExactOut].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	amountInMax, err := ParseBigInt(args["amountInMax"])
	if err != nil {
		return err
	}
	amountOut, err := ParseBigInt(args["amountOut"])
	if err != nil {
		return err
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return err
	}
	to, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	payerIsUser := args["payerIsUser"].(bool)
	s.AmountInMax = amountInMax
	s.AmountOut = amountOut
	s.Path = path
	s.To = to
	s.PayerIsUser = payerIsUser
	return nil
}

func (u *UniswapClient) V3SwapExactOut(tx MevTx, args map[string]interface{}) {}

func (s *JSONV3SwapExactOutParams) ConvertToBigIntType() *V3SwapExactOutParams {
	amountInMax, _ := new(big.Int).SetString(s.AmountInMax, 10)
	amountOut, _ := new(big.Int).SetString(s.AmountOut, 10)
	return &V3SwapExactOutParams{
		AmountInMax: amountInMax,
		AmountOut:   amountOut,
		Path:        s.Path,
		To:          s.To,
		PayerIsUser: s.PayerIsUser,
	}
}

func (s *V3SwapExactOutParams) ConvertToJSONType() *JSONV3SwapExactOutParams {
	return &JSONV3SwapExactOutParams{
		AmountInMax: s.AmountInMax.String(),
		AmountOut:   s.AmountOut.String(),
		Path:        s.Path,
		To:          s.To,
		PayerIsUser: s.PayerIsUser,
	}
}

func (s *V3SwapExactInParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		Trade: Trade{
			TradeMethod:             V3SwapExactIn,
			JSONV3SwapExactInParams: s.ConvertToJSONType(),
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

func (s *V3SwapExactOutParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountInMax)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		Trade: Trade{
			TradeMethod:              V3SwapExactOut,
			JSONV3SwapExactOutParams: s.ConvertToJSONType(),
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

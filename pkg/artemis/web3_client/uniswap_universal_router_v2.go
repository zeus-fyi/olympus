package web3_client

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
)

const (
	V2SwapExactIn  = "V2_SWAP_EXACT_IN"
	V2SwapExactOut = "V2_SWAP_EXACT_OUT"
)

type V2SwapExactInParams struct {
	AmountIn      *big.Int           `json:"amountIn"`
	AmountOutMin  *big.Int           `json:"amountOutMin"`
	Path          []accounts.Address `json:"path"`
	To            accounts.Address   `json:"to"`
	PayerIsSender bool               `json:"inputFromSender"`
}

type JSONV2SwapExactInParams struct {
	AmountIn      string             `json:"amountIn"`
	AmountOutMin  string             `json:"amountOutMin"`
	Path          []accounts.Address `json:"path"`
	To            accounts.Address   `json:"to"`
	PayerIsSender bool               `json:"payerIsSender"`
}

func (s *V2SwapExactInParams) Encode(ctx context.Context, abiFile *abi.ABI) ([]byte, error) {
	if abiFile == nil {
		inputs, err := UniversalRouterDecoderAbi.Methods[V2SwapExactIn].Inputs.Pack(s.To, s.AmountIn, s.AmountOutMin, s.Path, s.PayerIsSender)
		if err != nil {
			log.Err(err).Msg("V2SwapExactInParams: UniversalRouterDecoderAbi failed to encode")
			return nil, err
		}
		return inputs, nil
	} else {
		inputs, err := abiFile.Methods[V2SwapExactIn].Inputs.Pack(s.To, s.AmountIn, s.AmountOutMin, s.Path, s.PayerIsSender)
		if err != nil {
			log.Err(err).Msg("V2SwapExactInParams: abiFile failed to encode")
			return nil, err
		}
		return inputs, nil
	}
}

func (s *V2SwapExactInParams) Decode(ctx context.Context, data []byte, abiFile *abi.ABI) error {
	args := make(map[string]interface{})
	if abiFile == nil {
		err := UniversalRouterDecoderAbi.Methods[V2SwapExactIn].Inputs.UnpackIntoMap(args, data)
		if err != nil {
			log.Warn().Msg("V2SwapExactInParams: UniversalRouterDecoderAbi failed to decode")
			log.Err(err).Msg("V2SwapExactInParams: UniversalRouterDecoderAbi failed to decode")
			return err
		}
	} else {
		err := abiFile.Methods[V2SwapExactIn].Inputs.UnpackIntoMap(args, data)
		if err != nil {
			log.Warn().Msg("V2SwapExactInParams: abiFile failed to decode")
			log.Err(err).Msg("V2SwapExactInParams: abiFile failed to decode")
			return err
		}
	}
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		log.Warn().Msg("V2SwapExactInParams: failed to parse amountIn")
		log.Err(err).Msg("V2SwapExactInParams: failed to parse amountIn")
		return err
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		log.Warn().Msg("V2SwapExactInParams: failed to parse amountOutMin")
		log.Err(err).Msg("V2SwapExactInParams: failed to parse amountOutMin")
		return err
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		log.Warn().Msg("V2SwapExactInParams: failed to parse path")
		log.Err(err).Msg("V2SwapExactInParams: failed to parse path")
		return err
	}
	to, err := ConvertToAddress(args["recipient"])
	if err != nil {
		log.Warn().Msg("V2SwapExactInParams: failed to parse recipient")
		log.Err(err).Msg("V2SwapExactInParams: failed to parse recipient")
		return err
	}
	payerIsSender, ok := args["payerIsSender"].(bool)
	if !ok {
		log.Warn().Msg("V2SwapExactInParams: payerIsSender is not a bool, defaulting to false")
		payerIsSender = false
	}
	s.AmountIn = amountIn
	s.AmountOutMin = amountOutMin
	s.Path = path
	s.To = to
	s.PayerIsSender = payerIsSender
	return nil
}

func (s *JSONV2SwapExactInParams) ConvertToBigIntType() (*V2SwapExactInParams, error) {
	amountIn, ok1 := new(big.Int).SetString(s.AmountIn, 10)
	amountOutMin, ok2 := new(big.Int).SetString(s.AmountOutMin, 10)
	if !ok1 || !ok2 {
		log.Warn().Msg("JSONV2SwapExactInParams: failed to parse amountIn or amountOutMin")
		return nil, errors.New("JSONV2SwapExactInParams: failed to parse amountIn or amountOutMin")
	}
	return &V2SwapExactInParams{
		AmountIn:      amountIn,
		AmountOutMin:  amountOutMin,
		Path:          s.Path,
		To:            s.To,
		PayerIsSender: s.PayerIsSender,
	}, nil
}

func (s *V2SwapExactInParams) ConvertToJSONType() *JSONV2SwapExactInParams {
	return &JSONV2SwapExactInParams{
		AmountIn:      s.AmountIn.String(),
		AmountOutMin:  s.AmountOutMin.String(),
		Path:          s.Path,
		To:            s.To,
		PayerIsSender: s.PayerIsSender,
	}
}

type V2SwapExactOutParams struct {
	AmountInMax   *big.Int           `json:"amountInMax"`
	AmountOut     *big.Int           `json:"amountOut"`
	Path          []accounts.Address `json:"path"`
	To            accounts.Address   `json:"to"`
	PayerIsSender bool               `json:"payerIsSender"`
}

type JSONV2SwapExactOutParams struct {
	AmountInMax   string             `json:"amountInMax"`
	AmountOut     string             `json:"amountOut"`
	Path          []accounts.Address `json:"path"`
	To            accounts.Address   `json:"to"`
	PayerIsSender bool               `json:"payerIsSender"`
}

func (s *V2SwapExactOutParams) Encode(ctx context.Context, abiFile *abi.ABI) ([]byte, error) {
	if abiFile == nil {
		inputs, err := UniversalRouterDecoderAbi.Methods[V2SwapExactOut].Inputs.Pack(s.To, s.AmountOut, s.AmountInMax, s.Path, s.PayerIsSender)
		if err != nil {
			log.Error().Err(err).Msg("Failed to encode")
			return nil, err
		}
		return inputs, nil
	} else {
		inputs, err := abiFile.Methods[V2SwapExactOut].Inputs.Pack(s.To, s.AmountOut, s.AmountInMax, s.Path, s.PayerIsSender)
		if err != nil {
			log.Error().Err(err).Msg("Failed to encode")
			return nil, err
		}
		return inputs, nil
	}
}

func (s *V2SwapExactOutParams) Decode(ctx context.Context, data []byte, abiFile *abi.ABI) error {
	args := make(map[string]interface{})
	if abiFile == nil {
		err := UniversalRouterDecoderAbi.Methods[V2SwapExactOut].Inputs.UnpackIntoMap(args, data)
		if err != nil {
			log.Warn().Msg("V2SwapExactOutParams: UniversalRouterDecoderAbi failed to unpack")
			return err
		}
	} else {
		err := abiFile.Methods[V2SwapExactOut].Inputs.UnpackIntoMap(args, data)
		if err != nil {
			log.Warn().Msg("V2SwapExactOutParams: abiFile failed to unpack")
			return err
		}
	}
	amountInMax, err := ParseBigInt(args["amountInMax"])
	if err != nil {
		log.Err(err).Msg("V2SwapExactOutParams: failed to parse amountInMax")
		return err
	}
	amountOut, err := ParseBigInt(args["amountOut"])
	if err != nil {
		log.Err(err).Msg("V2SwapExactOutParams: failed to parse amountOut")
		return err
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		log.Err(err).Msg("V2SwapExactOutParams: failed to parse path")
		return err
	}
	to, err := ConvertToAddress(args["recipient"])
	if err != nil {
		log.Err(err).Msg("V2SwapExactOutParams: failed to parse recipient")
		return err
	}
	payerIsSender, ok := args["payerIsSender"].(bool)
	if !ok {
		log.Warn().Msg("V2SwapExactOutParams: payerIsSender is not a bool, defaulting to false")
		payerIsSender = false
	}
	s.AmountInMax = amountInMax
	s.AmountOut = amountOut
	s.Path = path
	s.To = to
	s.PayerIsSender = payerIsSender
	return nil
}

func (s *JSONV2SwapExactOutParams) ConvertToBigIntType() (*V2SwapExactOutParams, error) {
	amountInMax, ok := new(big.Int).SetString(s.AmountInMax, 10)
	if !ok {
		return nil, errors.New("failed to convert amountInMax to big.Int")
	}
	amountOut, ok := new(big.Int).SetString(s.AmountOut, 10)
	if !ok {
		return nil, errors.New("failed to convert amountOut to big.Int")
	}
	return &V2SwapExactOutParams{
		AmountInMax:   amountInMax,
		AmountOut:     amountOut,
		Path:          s.Path,
		To:            s.To,
		PayerIsSender: s.PayerIsSender,
	}, nil
}

func (s *V2SwapExactOutParams) ConvertToJSONType() *JSONV2SwapExactOutParams {
	return &JSONV2SwapExactOutParams{
		AmountInMax:   s.AmountInMax.String(),
		AmountOut:     s.AmountOut.String(),
		Path:          s.Path,
		To:            s.To,
		PayerIsSender: s.PayerIsSender,
	}
}

func (s *V2SwapExactOutParams) BinarySearch(pair uniswap_pricing.UniswapV2Pair) (TradeExecutionFlow, error) {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountInMax)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		InitialPair: &pair,
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

func (s *V2SwapExactInParams) BinarySearch(pair uniswap_pricing.UniswapV2Pair) (TradeExecutionFlow, error) {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		InitialPair: &pair,
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

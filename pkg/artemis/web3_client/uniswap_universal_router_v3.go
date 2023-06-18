package web3_client

import (
	"context"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

const (
	V3SwapExactIn  = "V3_SWAP_EXACT_IN"
	V3SwapExactOut = "V3_SWAP_EXACT_OUT"
)

type V3SwapExactInParams struct {
	AmountIn     *big.Int         `json:"amountIn"`
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         TokenFeePath     `json:"path"`
	To           accounts.Address `json:"to"`
	PayerIsUser  bool             `json:"payerIsUser"`
}

type JSONV3SwapExactInParams struct {
	AmountIn     string           `json:"amountIn"`
	AmountOutMin string           `json:"amountOutMin"`
	Path         TokenFeePath     `json:"path"`
	To           accounts.Address `json:"to"`
	PayerIsUser  bool             `json:"payerIsUser"`
}

func (s *V3SwapExactInParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[V3SwapExactIn].Inputs.Pack(s.To, s.AmountIn, s.AmountOutMin, s.Path.Encode(), s.PayerIsUser)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

type TokenFee struct {
	Token accounts.Address
	Fee   *big.Int
}

type TokenFeePath struct {
	TokenIn accounts.Address
	Path    []TokenFee
}

func (tfp *TokenFeePath) Encode() []byte {
	// Convert TokenIn into bytes
	tokenIn := tfp.TokenIn.Bytes()

	// Initialize a slice to hold the path bytes
	var pathBytes []byte

	// Iterate over the path to encode each TokenFee into bytes
	for _, tf := range tfp.Path {
		// Convert each TokenFee's token into bytes
		token := tf.Token.Bytes()

		// Convert each TokenFee's fee into a 3 bytes (6 hex characters)
		feeBytes := big.NewInt(tf.Fee.Int64()).Bytes()
		// If feeBytes is not 3 bytes long, pad it with leading zeros
		for len(feeBytes) < 3 {
			feeBytes = append([]byte{0}, feeBytes...)
		}

		// Append the fee and token bytes to the pathBytes
		pathBytes = append(pathBytes, feeBytes...)
		pathBytes = append(pathBytes, token...)
	}

	// Concatenate TokenIn and Path bytes
	return append(tokenIn, pathBytes...)
}

func (tfp *TokenFeePath) GetEndToken() accounts.Address {
	return tfp.Path[len(tfp.Path)-1].Token
}

func (tfp *TokenFeePath) GetPath() []accounts.Address {
	path := []accounts.Address{tfp.TokenIn}
	for _, p := range tfp.Path {
		path = append(path, p.Token)
	}
	return path
}

func (tfp *TokenFeePath) Reverse() {
	pathList := tfp.Path
	for i, j := 0, len(pathList)-1; i < j; i, j = i+1, j-1 {
		pathList[i], pathList[j] = pathList[j], pathList[i]
	}
	pathList[len(pathList)-1].Token, tfp.TokenIn = tfp.TokenIn, pathList[len(pathList)-1].Token
}
func (s *V3SwapExactInParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[V3SwapExactIn].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		return err
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		return err
	}
	pathBytes := args["path"].([]byte)
	hexStr := accounts.Bytes2Hex(pathBytes)
	tfp := TokenFeePath{
		TokenIn: accounts.HexToAddress(hexStr[:40]),
	}
	var pathList []TokenFee
	for i := 0; i < len(hexStr[40:]); i += 46 {
		fee, _ := new(big.Int).SetString(hexStr[40:][i:i+6], 16)
		token := accounts.HexToAddress(hexStr[40:][i+6 : i+46])
		tf := TokenFee{
			Token: token,
			Fee:   fee,
		}
		pathList = append(pathList, tf)
	}
	tfp.Path = pathList
	tfp.Reverse()

	to, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	payerIsSender := args["payerIsUser"].(bool)
	s.AmountIn = amountIn
	s.AmountOutMin = amountOutMin
	s.Path = tfp
	s.To = to
	s.PayerIsUser = payerIsSender
	return err
}

func (u *UniswapClient) V3SwapExactIn(pair UniswapV2Pair, inputs *V3SwapExactInParams) {
	tf := inputs.BinarySearch(pair)
	tf.InitialPair = pair.ConvertToJSONType()
}

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

type V3SwapExactOutParams struct {
	AmountInMax *big.Int         `json:"amountInMax"`
	AmountOut   *big.Int         `json:"amountOut"`
	Path        TokenFeePath     `json:"path"`
	To          accounts.Address `json:"to"`
	PayerIsUser bool             `json:"payerIsUser"`
}

type JSONV3SwapExactOutParams struct {
	AmountInMax string           `json:"amountInMax"`
	AmountOut   string           `json:"amountOut"`
	Path        TokenFeePath     `json:"path"`
	To          accounts.Address `json:"to"`
	PayerIsUser bool             `json:"payerIsUser"`
}

func (s *V3SwapExactOutParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[V3SwapExactOut].Inputs.Pack(s.To, s.AmountOut, s.AmountInMax, s.Path.Encode(), s.PayerIsUser)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (s *V3SwapExactOutParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[V3SwapExactOut].Inputs.UnpackIntoMap(args, data)
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
	pathBytes := args["path"].([]byte)
	hexStr := accounts.Bytes2Hex(pathBytes)
	tfp := TokenFeePath{
		TokenIn: accounts.HexToAddress(hexStr[:40]),
	}
	var pathList []TokenFee
	for i := 0; i < len(hexStr[40:]); i += 46 {
		fee, _ := new(big.Int).SetString(hexStr[40:][i:i+6], 16)
		token := accounts.HexToAddress(hexStr[40:][i+6 : i+46])
		tf := TokenFee{
			Token: token,
			Fee:   fee,
		}
		pathList = append(pathList, tf)
	}
	tfp.Path = pathList
	to, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	payerIsUser := args["payerIsUser"].(bool)
	s.AmountInMax = amountInMax
	s.AmountOut = amountOut
	s.Path = tfp
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

func (s *V3SwapExactInParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlowJSON {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlowJSON{
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
		toFrontRun, err := mockPairResp.PriceImpact(s.Path.TokenIn, mid)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		// User trade
		to, err := mockPairResp.PriceImpact(s.Path.TokenIn, s.AmountIn)
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
		toSandwich, err := mockPairResp.PriceImpact(s.Path.GetEndToken(), sandwichDump)
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

func (s *V3SwapExactOutParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlowJSON {
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
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		mid = new(big.Int).Add(low, high)
		mid = DivideByHalf(mid)
		// Front run trade
		toFrontRun, err := mockPairResp.PriceImpact(s.Path.TokenIn, mid)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		// User trade
		to, err := mockPairResp.PriceImpact(s.Path.TokenIn, s.AmountInMax)
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
		toSandwich, err := mockPairResp.PriceImpact(s.Path.GetEndToken(), sandwichDump)
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

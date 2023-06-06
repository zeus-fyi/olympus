package web3_client

import (
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

type SwapETHForExactTokensParams struct {
	AmountOut *big.Int           `json:"amountOut"`
	Path      []accounts.Address `json:"path"`
	To        accounts.Address   `json:"to"`
	Deadline  *big.Int           `json:"deadline"`
	Value     *big.Int           `json:"value"`
}

type JSONSwapETHForExactTokensParams struct {
	AmountOut string             `json:"amountOut"`
	Path      []accounts.Address `json:"path"`
	To        accounts.Address   `json:"to"`
	Deadline  string             `json:"deadline"`
	Value     string             `json:"value"`
}

func (s *SwapETHForExactTokensParams) ConvertToJSONType() *JSONSwapETHForExactTokensParams {
	return &JSONSwapETHForExactTokensParams{
		AmountOut: s.AmountOut.String(),
		Path:      s.Path,
		To:        s.To,
		Deadline:  s.Deadline.String(),
		Value:     s.Value.String(),
	}
}
func (s *SwapETHForExactTokensParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	// Value == variable
	// AmountOut == required for trade
	low := big.NewInt(0)
	high := new(big.Int).Set(s.Value)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		Trade: Trade{
			TradeMethod:                     "swapETHForExactTokens",
			JSONSwapETHForExactTokensParams: s.ConvertToJSONType(),
		},
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		mid = new(big.Int).Add(low, high)
		mid = mid.Div(mid, big.NewInt(2))
		// Front run trade
		toFrontRun, err := mockPairResp.PriceImpact(s.Path[0], mid)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		// User trade
		to, err := mockPairResp.PriceImpact(s.Path[0], s.Value)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOut)
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

type SwapTokensForExactTokensParams struct {
	AmountOut   *big.Int           `json:"amountOut"`
	AmountInMax *big.Int           `json:"amountInMax"`
	Path        []accounts.Address `json:"path"`
	To          accounts.Address   `json:"to"`
	Deadline    *big.Int           `json:"deadline"`
}

type JSONSwapTokensForExactTokensParams struct {
	AmountOut   string             `json:"amountOut"`
	AmountInMax string             `json:"amountInMax"`
	Path        []accounts.Address `json:"path"`
	To          accounts.Address   `json:"to"`
	Deadline    string             `json:"deadline"`
}

func (s *SwapTokensForExactTokensParams) ConvertToJSONType() *JSONSwapTokensForExactTokensParams {
	return &JSONSwapTokensForExactTokensParams{
		AmountOut:   s.AmountOut.String(),
		AmountInMax: s.AmountInMax.String(),
		Path:        s.Path,
		To:          s.To,
		Deadline:    s.Deadline.String(),
	}
}

func (s *SwapTokensForExactTokensParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountInMax)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		Trade: Trade{
			TradeMethod:                        "swapTokensForExactTokens",
			JSONSwapTokensForExactTokensParams: s.ConvertToJSONType(),
		},
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		mid = new(big.Int).Add(low, high)
		mid = mid.Div(mid, big.NewInt(2))
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

type SwapTokensForExactETHParams struct {
	AmountOut   *big.Int           `json:"amountOut"`
	AmountInMax *big.Int           `json:"amountInMax"`
	Path        []accounts.Address `json:"path"`
	To          accounts.Address   `json:"to"`
	Deadline    *big.Int           `json:"deadline"`
}

type JSONSwapTokensForExactETHParams struct {
	AmountOut   string             `json:"amountOut"`
	AmountInMax string             `json:"amountInMax"`
	Path        []accounts.Address `json:"path"`
	To          accounts.Address   `json:"to"`
	Deadline    string             `json:"deadline"`
}

func (s *SwapTokensForExactETHParams) ConvertToJSONType() *JSONSwapTokensForExactETHParams {
	return &JSONSwapTokensForExactETHParams{
		AmountOut:   s.AmountOut.String(),
		AmountInMax: s.AmountInMax.String(),
		Path:        s.Path,
		To:          s.To,
		Deadline:    s.Deadline.String(),
	}
}
func (s *SwapTokensForExactETHParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountInMax)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		Trade: Trade{
			TradeMethod:                     "swapTokensForExactETH",
			JSONSwapTokensForExactETHParams: s.ConvertToJSONType(),
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

func DivideByHalf(input *big.Int) *big.Int {
	modEven := new(big.Int).Mod(input, big.NewInt(2))
	if modEven.String() == "0" {
		input = input.Div(input, big.NewInt(2))
	} else {
		input = input.Add(input, big.NewInt(1))
		input = input.Div(input, big.NewInt(2))
	}
	return input
}

type SandwichTradePrediction struct {
	SellAmount     *big.Int `json:"sellAmount"`
	ExpectedProfit *big.Int `json:"expectedProfit"`
}

type JSONSandwichTradePrediction struct {
	SellAmount     string `json:"sellAmount"`
	ExpectedProfit string `json:"expectedProfit"`
}

func (s *JSONSandwichTradePrediction) ConvertToBigIntType() SandwichTradePrediction {
	sellAmount, _ := new(big.Int).SetString(s.SellAmount, 10)
	expectedProfit, _ := new(big.Int).SetString(s.ExpectedProfit, 10)
	return SandwichTradePrediction{
		SellAmount:     sellAmount,
		ExpectedProfit: expectedProfit,
	}
}

func (s *SandwichTradePrediction) ConvertToJSONType() JSONSandwichTradePrediction {
	if s.SellAmount == nil {
		s.SellAmount = big.NewInt(0)
	}
	if s.ExpectedProfit == nil {
		s.ExpectedProfit = big.NewInt(0)
	}
	return JSONSandwichTradePrediction{
		SellAmount:     s.SellAmount.String(),
		ExpectedProfit: s.ExpectedProfit.String(),
	}
}

type SwapExactTokensForTokensParams struct {
	AmountIn     *big.Int           `json:"amountIn"`
	AmountOutMin *big.Int           `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     *big.Int           `json:"deadline"`
}

type JSONSwapExactTokensForTokensParams struct {
	AmountIn     string             `json:"amountIn"`
	AmountOutMin string             `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     string             `json:"deadline"`
}

func (s *SwapExactTokensForTokensParams) ConvertToJSONType() *JSONSwapExactTokensForTokensParams {
	return &JSONSwapExactTokensForTokensParams{
		AmountIn:     s.AmountIn.String(),
		AmountOutMin: s.AmountOutMin.String(),
		Path:         s.Path,
		To:           s.To,
		Deadline:     s.Deadline.String(),
	}
}

func (s *SwapExactTokensForTokensParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		Trade: Trade{
			TradeMethod:                        "swapExactTokensForTokens",
			JSONSwapExactTokensForTokensParams: s.ConvertToJSONType(),
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
func (s *SwapExactETHForTokensParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.Value)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		Trade: Trade{
			TradeMethod:                     "swapExactETHForTokens",
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
			return tf
		}
		// User trade
		to, err := mockPairResp.PriceImpact(s.Path[0], s.Value)
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

type SwapExactTokensForETHParams struct {
	AmountIn     *big.Int           `json:"amountIn"`
	AmountOutMin *big.Int           `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     *big.Int           `json:"deadline"`
}

type JSONSwapExactTokensForETHParams struct {
	AmountIn     string             `json:"amountIn"`
	AmountOutMin string             `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     string             `json:"deadline"`
}

func (s *SwapExactTokensForETHParams) ConvertToJSONType() *JSONSwapExactTokensForETHParams {
	return &JSONSwapExactTokensForETHParams{
		AmountIn:     s.AmountIn.String(),
		AmountOutMin: s.AmountOutMin.String(),
		Path:         s.Path,
		To:           s.To,
		Deadline:     s.Deadline.String(),
	}
}
func (s *SwapExactTokensForETHParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		Trade: Trade{
			TradeMethod:                     "swapExactTokensForETH",
			JSONSwapExactTokensForETHParams: s.ConvertToJSONType(),
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

type SwapExactTokensForTokensSupportingFeeOnTransferTokensParams struct {
	AmountIn     *big.Int           `json:"amountIn"`
	AmountOutMin *big.Int           `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     *big.Int           `json:"deadline"`
}

type SwapExactETHForTokensSupportingFeeOnTransferTokensParams struct {
	AmountOutMin *big.Int           `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     *big.Int           `json:"deadline"`
}

type SwapExactTokensForETHSupportingFeeOnTransferTokensParams struct {
	AmountIn     *big.Int           `json:"amountIn"`
	AmountOutMin *big.Int           `json:"amountOutMin"`
	Path         []accounts.Address `json:"path"`
	To           accounts.Address   `json:"to"`
	Deadline     *big.Int           `json:"deadline"`
}

type AddLiquidityParams struct {
	TokenA         accounts.Address `json:"tokenA"`
	TokenB         accounts.Address `json:"tokenB"`
	AmountADesired *big.Int         `json:"amountADesired"`
	AmountBDesired *big.Int         `json:"amountBDesired"`
	AmountAMin     *big.Int         `json:"amountAMin"`
	AmountBMin     *big.Int         `json:"amountBMin"`
	To             accounts.Address `json:"to"`
	Deadline       *big.Int         `json:"deadline"`
}

type AddLiquidityETHParams struct {
	Token              accounts.Address `json:"token"`
	AmountTokenDesired *big.Int         `json:"amountTokenDesired"`
	AmountTokenMin     *big.Int         `json:"amountTokenMin"`
	AmountETHMin       *big.Int         `json:"amountETHMin"`
	To                 accounts.Address `json:"to"`
	Deadline           *big.Int         `json:"deadline"`
}

type RemoveLiquidityParams struct {
	TokenA     accounts.Address `json:"tokenA"`
	TokenB     accounts.Address `json:"tokenB"`
	Liquidity  *big.Int         `json:"liquidity"`
	AmountAMin *big.Int         `json:"amountAMin"`
	AmountBMin *big.Int         `json:"amountBMin"`
	To         accounts.Address `json:"to"`
	Deadline   *big.Int         `json:"deadline"`
}

type RemoveLiquidityETHParams struct {
	Token          accounts.Address `json:"token"`
	Liquidity      *big.Int         `json:"liquidity"`
	AmountTokenMin *big.Int         `json:"amountTokenMin"`
	AmountETHMin   *big.Int         `json:"amountETHMin"`
	To             accounts.Address `json:"to"`
	Deadline       *big.Int         `json:"deadline"`
}

type RemoveLiquidityWithPermitParams struct {
	TokenA     accounts.Address `json:"tokenA"`
	TokenB     accounts.Address `json:"tokenB"`
	Liquidity  *big.Int         `json:"liquidity"`
	AmountAMin *big.Int         `json:"amountAMin"`
	AmountBMin *big.Int         `json:"amountBMin"`
	To         accounts.Address `json:"to"`
	Deadline   *big.Int         `json:"deadline"`
	ApproveMax bool             `json:"approveMax"`
	V          uint8            `json:"v"`
	R          [32]byte         `json:"r"`
	S          [32]byte         `json:"s"`
}

type RemoveLiquidityETHWithPermitParams struct {
	Token          accounts.Address `json:"token"`
	Liquidity      *big.Int         `json:"liquidity"`
	AmountTokenMin *big.Int         `json:"amountTokenMin"`
	AmountETHMin   *big.Int         `json:"amountETHMin"`
	To             accounts.Address `json:"to"`
	Deadline       *big.Int         `json:"deadline"`
	ApproveMax     bool             `json:"approveMax"`
	V              uint8            `json:"v"`
	R              [32]byte         `json:"r"`
	S              [32]byte         `json:"s"`
}

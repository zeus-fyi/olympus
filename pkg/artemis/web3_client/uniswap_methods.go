package web3_client

import (
	"fmt"
	"math/big"

	"github.com/gochain/gochain/v4/common"
	web3_types "github.com/zeus-fyi/gochain/web3/types"
)

type TradeExecutionFlow struct {
	CurrentBlockNumber *big.Int                   `json:"currentBlockNumber"`
	Tx                 *web3_types.RpcTransaction `json:"tx"`
	TradeMethod        string                     `json:"tradeMethod"`
	TradeParams        any                        `json:"tradeParams"`
	InitialPair        UniswapV2Pair              `json:"initialPair"`
	FrontRunTrade      TradeOutcome               `json:"frontRunTrade"`
	UserTrade          TradeOutcome               `json:"userTrade"`
	SandwichTrade      TradeOutcome               `json:"sandwichTrade"`
	SandwichPrediction SandwichTradePrediction    `json:"sandwichPrediction"`
}

type SwapETHForExactTokensParams struct {
	AmountOut *big.Int         `json:"amountOut"`
	Path      []common.Address `json:"path"`
	To        common.Address   `json:"to"`
	Deadline  *big.Int         `json:"deadline"`
	Value     *big.Int         `json:"value"`
}

func (s *SwapETHForExactTokensParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	// Value == variable
	// AmountOut == required for trade
	low := big.NewInt(0)
	high := new(big.Int).Set(s.Value)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmount *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		TradeMethod: "swapExactTokensForTokens",
		TradeParams: s,
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		tokenSellAmount = new(big.Int).Add(low, high)
		mid = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
		// Front run trade
		toFrontRun := mockPairResp.PriceImpact(s.Path[0], tokenSellAmount)
		// User trade
		to := mockPairResp.PriceImpact(s.Path[0], s.Value)
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOut)
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
			continue
		}
		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich := mockPairResp.PriceImpact(s.Path[1], sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			mid = tokenSellAmount
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = toFrontRun
			tf.UserTrade = to
			tf.SandwichTrade = toSandwich
		}
		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	sp := SandwichTradePrediction{
		SellAmount:     tokenSellAmountAtMaxProfit,
		ExpectedProfit: maxProfit,
	}
	tf.SandwichPrediction = sp
	return tf
}

type SwapTokensForExactTokensParams struct {
	AmountOut   *big.Int         `json:"amountOut"`
	AmountInMax *big.Int         `json:"amountInMax"`
	Path        []common.Address `json:"path"`
	To          common.Address   `json:"to"`
	Deadline    *big.Int         `json:"deadline"`
}

func (s *SwapTokensForExactTokensParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountInMax)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmount *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		TradeMethod: "swapTokensForExactETH",
		TradeParams: s,
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		tokenSellAmount = new(big.Int).Add(low, high)
		mid = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
		// Front run trade
		toFrontRun := mockPairResp.PriceImpact(s.Path[0], tokenSellAmount)
		// User trade
		to := mockPairResp.PriceImpact(s.Path[0], s.AmountInMax)
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOut)
		// if diff <= 0 then it searches left
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
			continue
		}
		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich := mockPairResp.PriceImpact(s.Path[1], sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			mid = tokenSellAmount
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = toFrontRun
			tf.UserTrade = to
			tf.SandwichTrade = toSandwich
		}
		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	sp := SandwichTradePrediction{
		SellAmount:     tokenSellAmountAtMaxProfit,
		ExpectedProfit: maxProfit,
	}
	tf.SandwichPrediction = sp
	return tf
}

type SwapTokensForExactETHParams struct {
	AmountOut   *big.Int         `json:"amountOut"`
	AmountInMax *big.Int         `json:"amountInMax"`
	Path        []common.Address `json:"path"`
	To          common.Address   `json:"to"`
	Deadline    *big.Int         `json:"deadline"`
}

func (s *SwapTokensForExactETHParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountInMax)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmount *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		TradeMethod: "swapTokensForExactETH",
		TradeParams: s,
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		tokenSellAmount = new(big.Int).Add(low, high)
		mid = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
		// Front run trade
		toFrontRun := mockPairResp.PriceImpact(s.Path[0], tokenSellAmount)
		// User trade
		to := mockPairResp.PriceImpact(s.Path[0], s.AmountInMax)
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOut)
		// if diff <= 0 then it searches left
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
			continue
		}
		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich := mockPairResp.PriceImpact(s.Path[1], sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			mid = tokenSellAmount
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = toFrontRun
			tf.UserTrade = to
			tf.SandwichTrade = toSandwich
		}
		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	sp := SandwichTradePrediction{
		SellAmount:     tokenSellAmountAtMaxProfit,
		ExpectedProfit: maxProfit,
	}
	tf.SandwichPrediction = sp
	return tf
}

type SandwichTradePrediction struct {
	SellAmount     *big.Int `json:"sellAmount"`
	ExpectedProfit *big.Int `json:"expectedProfit"`
}

type SwapExactTokensForTokensParams struct {
	AmountIn     *big.Int         `json:"amountIn"`
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         []common.Address `json:"path"`
	To           common.Address   `json:"to"`
	Deadline     *big.Int         `json:"deadline"`
}

func (s *SwapExactTokensForTokensParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmount *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		TradeMethod: "swapExactTokensForTokens",
		TradeParams: s,
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		tokenSellAmount = new(big.Int).Add(low, high)
		mid = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
		// Front run trade
		toFrontRun := mockPairResp.PriceImpact(s.Path[0], tokenSellAmount)
		// User trade
		to := mockPairResp.PriceImpact(s.Path[0], s.AmountIn)
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOutMin)
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
			continue
		}
		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich := mockPairResp.PriceImpact(s.Path[1], sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			mid = tokenSellAmount
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = toFrontRun
			tf.UserTrade = to
			tf.SandwichTrade = toSandwich
		}
		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	sp := SandwichTradePrediction{
		SellAmount:     tokenSellAmountAtMaxProfit,
		ExpectedProfit: maxProfit,
	}
	tf.SandwichPrediction = sp
	return tf
}

type SwapExactETHForTokensParams struct {
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         []common.Address `json:"path"`
	To           common.Address   `json:"to"`
	Value        *big.Int         `json:"value"`
	Deadline     *big.Int         `json:"deadline"`
}

func (s *SwapExactETHForTokensParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.Value)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmount *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		TradeMethod: "swapExactETHForTokens",
		TradeParams: s,
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		tokenSellAmount = new(big.Int).Add(low, high)
		mid = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
		// Front run trade
		toFrontRun := mockPairResp.PriceImpact(s.Path[0], tokenSellAmount)
		// User trade
		to := mockPairResp.PriceImpact(s.Path[0], s.Value)
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOutMin)
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
			continue
		}
		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich := mockPairResp.PriceImpact(s.Path[1], sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			mid = tokenSellAmount
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = toFrontRun
			tf.UserTrade = to
			tf.SandwichTrade = toSandwich
		}
		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	sp := SandwichTradePrediction{
		SellAmount:     tokenSellAmountAtMaxProfit,
		ExpectedProfit: maxProfit,
	}
	tf.SandwichPrediction = sp
	return tf
}

type SwapExactTokensForETHParams struct {
	AmountIn     *big.Int         `json:"amountIn"`
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         []common.Address `json:"path"`
	To           common.Address   `json:"to"`
	Deadline     *big.Int         `json:"deadline"`
}

func (s *SwapExactTokensForETHParams) BinarySearch(pair UniswapV2Pair) TradeExecutionFlow {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmount *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		TradeMethod: "swapExactTokensForETH",
		TradeParams: s,
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		tokenSellAmount = new(big.Int).Add(low, high)
		mid = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
		// Front run trade
		toFrontRun := mockPairResp.PriceImpact(s.Path[0], tokenSellAmount)
		// User trade
		to := mockPairResp.PriceImpact(s.Path[0], s.AmountIn)
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOutMin)
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
			continue
		}
		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich := mockPairResp.PriceImpact(s.Path[1], sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			mid = tokenSellAmount
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = toFrontRun
			tf.UserTrade = to
			tf.SandwichTrade = toSandwich
		}
		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	sp := SandwichTradePrediction{
		SellAmount:     tokenSellAmountAtMaxProfit,
		ExpectedProfit: maxProfit,
	}
	tf.SandwichPrediction = sp
	return tf
}

type SwapExactTokensForTokensSupportingFeeOnTransferTokensParams struct {
	AmountIn     *big.Int         `json:"amountIn"`
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         []common.Address `json:"path"`
	To           common.Address   `json:"to"`
	Deadline     *big.Int         `json:"deadline"`
}

type SwapExactETHForTokensSupportingFeeOnTransferTokensParams struct {
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         []common.Address `json:"path"`
	To           common.Address   `json:"to"`
	Deadline     *big.Int         `json:"deadline"`
}

type SwapExactTokensForETHSupportingFeeOnTransferTokensParams struct {
	AmountIn     *big.Int         `json:"amountIn"`
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         []common.Address `json:"path"`
	To           common.Address   `json:"to"`
	Deadline     *big.Int         `json:"deadline"`
}

type AddLiquidityParams struct {
	TokenA         common.Address `json:"tokenA"`
	TokenB         common.Address `json:"tokenB"`
	AmountADesired *big.Int       `json:"amountADesired"`
	AmountBDesired *big.Int       `json:"amountBDesired"`
	AmountAMin     *big.Int       `json:"amountAMin"`
	AmountBMin     *big.Int       `json:"amountBMin"`
	To             common.Address `json:"to"`
	Deadline       *big.Int       `json:"deadline"`
}

type AddLiquidityETHParams struct {
	Token              common.Address `json:"token"`
	AmountTokenDesired *big.Int       `json:"amountTokenDesired"`
	AmountTokenMin     *big.Int       `json:"amountTokenMin"`
	AmountETHMin       *big.Int       `json:"amountETHMin"`
	To                 common.Address `json:"to"`
	Deadline           *big.Int       `json:"deadline"`
}

type RemoveLiquidityParams struct {
	TokenA     common.Address `json:"tokenA"`
	TokenB     common.Address `json:"tokenB"`
	Liquidity  *big.Int       `json:"liquidity"`
	AmountAMin *big.Int       `json:"amountAMin"`
	AmountBMin *big.Int       `json:"amountBMin"`
	To         common.Address `json:"to"`
	Deadline   *big.Int       `json:"deadline"`
}

type RemoveLiquidityETHParams struct {
	Token          common.Address `json:"token"`
	Liquidity      *big.Int       `json:"liquidity"`
	AmountTokenMin *big.Int       `json:"amountTokenMin"`
	AmountETHMin   *big.Int       `json:"amountETHMin"`
	To             common.Address `json:"to"`
	Deadline       *big.Int       `json:"deadline"`
}

type RemoveLiquidityWithPermitParams struct {
	TokenA     common.Address `json:"tokenA"`
	TokenB     common.Address `json:"tokenB"`
	Liquidity  *big.Int       `json:"liquidity"`
	AmountAMin *big.Int       `json:"amountAMin"`
	AmountBMin *big.Int       `json:"amountBMin"`
	To         common.Address `json:"to"`
	Deadline   *big.Int       `json:"deadline"`
	ApproveMax bool           `json:"approveMax"`
	V          uint8          `json:"v"`
	R          [32]byte       `json:"r"`
	S          [32]byte       `json:"s"`
}

type RemoveLiquidityETHWithPermitParams struct {
	Token          common.Address `json:"token"`
	Liquidity      *big.Int       `json:"liquidity"`
	AmountTokenMin *big.Int       `json:"amountTokenMin"`
	AmountETHMin   *big.Int       `json:"amountETHMin"`
	To             common.Address `json:"to"`
	Deadline       *big.Int       `json:"deadline"`
	ApproveMax     bool           `json:"approveMax"`
	V              uint8          `json:"v"`
	R              [32]byte       `json:"r"`
	S              [32]byte       `json:"s"`
}

func ParseBigInt(i interface{}) (*big.Int, error) {
	switch v := i.(type) {
	case *big.Int:
		return i.(*big.Int), nil
	case string:
		base := 10
		result := new(big.Int)
		_, ok := result.SetString(v, base)
		if !ok {
			return nil, fmt.Errorf("failed to parse string '%s' into big.Int", v)
		}
		return result, nil
	case uint32:
		return big.NewInt(int64(v)), nil
	case int64:
		return big.NewInt(v), nil
	default:
		return nil, fmt.Errorf("input is not a string or int64")
	}
}

func ConvertToAddressSlice(i interface{}) ([]common.Address, error) {
	switch v := i.(type) {
	case []common.Address:
		return i.([]common.Address), nil
	default:
		fmt.Println(v)
		return nil, fmt.Errorf("input is not a []common.Address")
	}
}

func ConvertToAddress(i interface{}) (common.Address, error) {
	switch v := i.(type) {
	case common.Address:
		return i.(common.Address), nil
	default:
		fmt.Println(v)
		return common.Address{}, fmt.Errorf("input is not a  common.Address")
	}
}

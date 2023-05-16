package web3_client

import (
	"fmt"
	"math/big"

	"github.com/gochain/gochain/v4/common"
)

type SwapExactTokensForTokensParams struct {
	AmountIn     *big.Int         `json:"amountIn"`
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         []common.Address `json:"path"`
	To           common.Address   `json:"to"`
	Deadline     *big.Int         `json:"deadline"`
}

type SwapTokensForExactTokensParams struct {
	AmountOut   *big.Int         `json:"amountOut"`
	AmountInMax *big.Int         `json:"amountInMax"`
	Path        []common.Address `json:"path"`
	To          common.Address   `json:"to"`
	Deadline    *big.Int         `json:"deadline"`
}

type SwapExactETHForTokensParams struct {
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         []common.Address `json:"path"`
	To           common.Address   `json:"to"`
	Value        *big.Int         `json:"value"`
	Deadline     *big.Int         `json:"deadline"`
}

type SwapTokensForExactETHParams struct {
	AmountOut   *big.Int         `json:"amountOut"`
	AmountInMax *big.Int         `json:"amountInMax"`
	Path        []common.Address `json:"path"`
	To          common.Address   `json:"to"`
	Deadline    *big.Int         `json:"deadline"`
}

type SwapExactTokensForETHParams struct {
	AmountIn     *big.Int         `json:"amountIn"`
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         []common.Address `json:"path"`
	To           common.Address   `json:"to"`
	Deadline     *big.Int         `json:"deadline"`
}

type SandwichTrade struct {
	SellAmount     *big.Int `json:"sellAmount"`
	ExpectedProfit *big.Int `json:"expectedProfit"`
}

func (s *SwapExactTokensForTokensParams) BinarySearch(pair UniswapV2Pair) SandwichTrade {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmount *big.Int
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		tokenSellAmount = new(big.Int).Add(low, high)
		mid = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
		// Front run trade
		toFrontRun, _, _ := mockPairResp.PriceImpactToken0BuyToken1(tokenSellAmount)

		// User trade
		to, _, _ := mockPairResp.PriceImpactToken0BuyToken1(s.AmountIn)
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOutMin)
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
			continue
		}

		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich, _, _ := mockPairResp.PriceImpactToken1BuyToken0(sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			mid = tokenSellAmount
		}

		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	return SandwichTrade{
		SellAmount:     mid,
		ExpectedProfit: maxProfit,
	}
}

func (s *SwapExactETHForTokensParams) BinarySearch(pair UniswapV2Pair) SandwichTrade {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.Value)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmount *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		tokenSellAmount = new(big.Int).Add(low, high)
		mid = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
		// Front run trade
		toFrontRun, _, _ := mockPairResp.PriceImpactToken0BuyToken1(tokenSellAmount)

		// User trade
		to, _, _ := mockPairResp.PriceImpactToken0BuyToken1(s.Value)
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOutMin)
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
			continue
		}

		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich, _, _ := mockPairResp.PriceImpactToken1BuyToken0(sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			mid = tokenSellAmount
			tokenSellAmountAtMaxProfit = mid
		}

		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	return SandwichTrade{
		SellAmount:     tokenSellAmountAtMaxProfit,
		ExpectedProfit: maxProfit,
	}
}

// SwapExactTokensForETH TODO verify
func (s *SwapExactTokensForETHParams) BinarySearch(pair UniswapV2Pair) SandwichTrade {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmount *big.Int
	var tokenSellAmountAtMaxProfit *big.Int

	for low.Cmp(high) <= 0 {
		mockPairResp := pair
		tokenSellAmount = new(big.Int).Add(low, high)
		mid = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
		// Front run trade
		toFrontRun, _, _ := mockPairResp.PriceImpactToken1BuyToken0(tokenSellAmount)

		// User trade
		to, _, _ := mockPairResp.PriceImpactToken1BuyToken0(s.AmountIn)
		difference := new(big.Int).Sub(to.AmountOut, s.AmountOutMin)
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
			continue
		}

		// Sandwich trade
		sandwichDump := toFrontRun.AmountOut
		toSandwich, _, _ := mockPairResp.PriceImpactToken0BuyToken1(sandwichDump)
		profit := new(big.Int).Sub(toSandwich.AmountOut, toFrontRun.AmountIn)
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			mid = tokenSellAmount
			tokenSellAmountAtMaxProfit = mid
		}

		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	return SandwichTrade{
		SellAmount:     tokenSellAmountAtMaxProfit,
		ExpectedProfit: maxProfit,
	}
}

type SwapETHForExactTokensParams struct {
	AmountOut *big.Int         `json:"amountOut"`
	Path      []common.Address `json:"path"`
	To        common.Address   `json:"to"`
	Deadline  *big.Int         `json:"deadline"`
	Value     *big.Int         `json:"value"`
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

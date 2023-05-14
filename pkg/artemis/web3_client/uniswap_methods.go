package web3_client

import (
	"fmt"
	"math/big"

	"github.com/gochain/gochain/v4/common"
)

type SwapExactTokensForTokensParams struct {
	AmountIn     *big.Int
	AmountOutMin *big.Int
	Path         []common.Address
	To           common.Address
	Deadline     *big.Int
}

type SwapTokensForExactTokensParams struct {
	AmountOut   *big.Int
	AmountInMax *big.Int
	Path        []common.Address
	To          common.Address
	Deadline    *big.Int
}

type SwapExactETHForTokensParams struct {
	AmountOutMin *big.Int
	Path         []common.Address
	To           common.Address
	Value        *big.Int
	Deadline     *big.Int
}

type SwapTokensForExactETHParams struct {
	AmountOut   *big.Int
	AmountInMax *big.Int
	Path        []common.Address
	To          common.Address
	Deadline    *big.Int
}

type SwapExactTokensForETHParams struct {
	AmountIn     *big.Int
	AmountOutMin *big.Int
	Path         []common.Address
	To           common.Address
	Deadline     *big.Int
}

func (s *SwapExactTokensForTokensParams) BinarySearch(pair UniswapV2Pair) (*big.Int, *big.Int) {
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
	fmt.Println("mid:", mid.String())
	return tokenSellAmount, maxProfit
}

func (s *SwapExactETHForTokensParams) BinarySearch(pair UniswapV2Pair) (*big.Int, *big.Int) {
	low := big.NewInt(0)
	high := new(big.Int).Set(s.Value)
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
		}

		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	fmt.Println("mid:", mid.String())
	return tokenSellAmount, maxProfit
}

func (s *SwapExactTokensForETHParams) BinarySearch(pair UniswapV2Pair) (*big.Int, *big.Int) {
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
		}

		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(tokenSellAmount, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(tokenSellAmount, big.NewInt(1))
		}
	}
	fmt.Println("mid:", mid.String())
	return tokenSellAmount, maxProfit
}

type SwapETHForExactTokensParams struct {
	AmountOut *big.Int
	Path      []common.Address
	To        common.Address
	Deadline  *big.Int
	Value     *big.Int
}

type SwapExactTokensForTokensSupportingFeeOnTransferTokensParams struct {
	AmountIn     *big.Int
	AmountOutMin *big.Int
	Path         []common.Address
	To           common.Address
	Deadline     *big.Int
}

type SwapExactETHForTokensSupportingFeeOnTransferTokensParams struct {
	AmountOutMin *big.Int
	Path         []common.Address
	To           common.Address
	Deadline     *big.Int
}

type SwapExactTokensForETHSupportingFeeOnTransferTokensParams struct {
	AmountIn     *big.Int
	AmountOutMin *big.Int
	Path         []common.Address
	To           common.Address
	Deadline     *big.Int
}

type AddLiquidityParams struct {
	TokenA         common.Address
	TokenB         common.Address
	AmountADesired *big.Int
	AmountBDesired *big.Int
	AmountAMin     *big.Int
	AmountBMin     *big.Int
	To             common.Address
	Deadline       *big.Int
}

type AddLiquidityETHParams struct {
	Token              common.Address
	AmountTokenDesired *big.Int
	AmountTokenMin     *big.Int
	AmountETHMin       *big.Int
	To                 common.Address
	Deadline           *big.Int
}

type RemoveLiquidityParams struct {
	TokenA     common.Address
	TokenB     common.Address
	Liquidity  *big.Int
	AmountAMin *big.Int
	AmountBMin *big.Int
	To         common.Address
	Deadline   *big.Int
}

type RemoveLiquidityETHParams struct {
	Token          common.Address
	Liquidity      *big.Int
	AmountTokenMin *big.Int
	AmountETHMin   *big.Int
	To             common.Address
	Deadline       *big.Int
}

type RemoveLiquidityWithPermitParams struct {
	TokenA     common.Address
	TokenB     common.Address
	Liquidity  *big.Int
	AmountAMin *big.Int
	AmountBMin *big.Int
	To         common.Address
	Deadline   *big.Int
	ApproveMax bool
	V          uint8
	R          [32]byte
	S          [32]byte
}

type RemoveLiquidityETHWithPermitParams struct {
	Token          common.Address
	Liquidity      *big.Int
	AmountTokenMin *big.Int
	AmountETHMin   *big.Int
	To             common.Address
	Deadline       *big.Int
	ApproveMax     bool
	V              uint8
	R              [32]byte
	S              [32]byte
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

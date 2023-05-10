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

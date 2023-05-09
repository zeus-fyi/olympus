package web3_client

import (
	"fmt"
	"math/big"

	"github.com/gochain/gochain/v4/common"
	"github.com/rs/zerolog/log"
)

type SwapExactTokensForTokensParams struct {
	AmountIn     *big.Int
	AmountOutMin *big.Int
	Path         []common.Address
	To           common.Address
	Deadline     uint64
	Amounts      []*big.Int
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
	case string:
		base := 10
		result := new(big.Int)
		_, ok := result.SetString(v, base)
		if !ok {
			return nil, fmt.Errorf("failed to parse string '%s' into big.Int", v)
		}
		return result, nil
	case int64:
		return big.NewInt(v), nil
	default:
		return nil, fmt.Errorf("input is not a string or int64")
	}
}

func ConvertToAddressSlice(input interface{}) ([]common.Address, error) {
	// First, we need to check that the input is actually a slice of strings.
	inputSlice, ok := input.([]interface{})
	if !ok {
		log.Info().Msgf("input is not a slice: %v", input)
		return nil, fmt.Errorf("input is not a slice")
	}
	addresses := make([]common.Address, len(inputSlice))
	for i, v := range inputSlice {
		addr := common.HexToAddress(v.(string))
		addresses[i] = addr
	}

	return addresses, nil
}

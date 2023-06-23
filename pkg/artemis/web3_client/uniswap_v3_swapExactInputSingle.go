package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

type SwapExactInputSingleArgs struct {
	TokenIn           accounts.Address `json:"tokenIn"`
	AmountIn          *big.Int         `json:"amountIn"`
	TokenOut          accounts.Address `json:"tokenOut"`
	AmountOutMinimum  *big.Int         `json:"amountOutMinimum"`
	Fee               *big.Int         `json:"fee"`
	Recipient         accounts.Address `json:"recipient"`
	SqrtPriceLimitX96 *big.Int         `json:"sqrtPriceLimitX96"`
}

type JSONSwapExactInputSingleArgs struct {
	TokenIn           accounts.Address `json:"tokenIn"`
	AmountIn          string           `json:"amountIn"`
	TokenOut          accounts.Address `json:"tokenOut"`
	AmountOutMinimum  string           `json:"amountOutMinimum"`
	Fee               string           `json:"fee"`
	Recipient         accounts.Address `json:"recipient"`
	SqrtPriceLimitX96 string           `json:"sqrtPriceLimitX96"`
}

func (s *JSONSwapExactInputSingleArgs) ConvertToBigIntType() SwapExactInputSingleArgs {
	amountIn, _ := new(big.Int).SetString(s.AmountIn, 10)
	amountOutMinimum, _ := new(big.Int).SetString(s.AmountOutMinimum, 10)
	fee, _ := new(big.Int).SetString(s.Fee, 10)
	sqrtPriceLimitX96, _ := new(big.Int).SetString(s.SqrtPriceLimitX96, 10)

	return SwapExactInputSingleArgs{
		TokenIn:           s.TokenIn,
		AmountIn:          amountIn,
		TokenOut:          s.TokenOut,
		AmountOutMinimum:  amountOutMinimum,
		Fee:               fee,
		Recipient:         s.Recipient,
		SqrtPriceLimitX96: sqrtPriceLimitX96,
	}
}

func (s *SwapExactInputSingleArgs) ConvertToJSONType() *JSONSwapExactInputSingleArgs {
	return &JSONSwapExactInputSingleArgs{
		TokenIn:           s.TokenIn,
		AmountIn:          s.AmountIn.String(),
		TokenOut:          s.TokenOut,
		AmountOutMinimum:  s.AmountOutMinimum.String(),
		Fee:               s.Fee.String(),
		Recipient:         s.Recipient,
		SqrtPriceLimitX96: s.SqrtPriceLimitX96.String(),
	}
}

func (s *SwapExactInputSingleArgs) Decode(ctx context.Context, args map[string]interface{}) error {
	params, ok := args["params"].(struct {
		TokenIn           common.Address "json:\"tokenIn\""
		TokenOut          common.Address "json:\"tokenOut\""
		Fee               *big.Int       "json:\"fee\""
		Recipient         common.Address "json:\"recipient\""
		AmountIn          *big.Int       "json:\"amountIn\""
		AmountOutMinimum  *big.Int       "json:\"amountOutMinimum\""
		SqrtPriceLimitX96 *big.Int       "json:\"sqrtPriceLimitX96\""
	})
	if !ok {
		return fmt.Errorf("params is not of the expected type")
	}
	s.TokenIn = accounts.Address(params.TokenIn)
	s.TokenOut = accounts.Address(params.TokenOut)
	s.Fee = params.Fee
	s.Recipient = accounts.Address(params.Recipient)
	s.AmountIn = params.AmountIn
	s.AmountOutMinimum = params.AmountOutMinimum
	s.SqrtPriceLimitX96 = params.SqrtPriceLimitX96
	return nil
}

package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

type SwapExactOutputSingleArgs struct {
	TokenIn           accounts.Address `json:"tokenIn"`
	AmountInMaximum   *big.Int         `json:"amountInMaximum"`
	TokenOut          accounts.Address `json:"tokenOut"`
	AmountOut         *big.Int         `json:"amountOut"`
	Fee               *big.Int         `json:"fee"`
	Recipient         accounts.Address `json:"recipient"`
	SqrtPriceLimitX96 *big.Int         `json:"sqrtPriceLimitX96"`
}

type JSONSwapExactOutputSingleArgs struct {
	TokenIn           accounts.Address `json:"tokenIn"`
	AmountOut         string           `json:"amountOut"`
	TokenOut          accounts.Address `json:"tokenOut"`
	AmountInMaximum   string           `json:"amountInMaximum"`
	Fee               string           `json:"fee"`
	Recipient         accounts.Address `json:"recipient"`
	SqrtPriceLimitX96 string           `json:"sqrtPriceLimitX96"`
}

func (s *JSONSwapExactOutputSingleArgs) ConvertToBigIntType() SwapExactOutputSingleArgs {
	amountOut, _ := new(big.Int).SetString(s.AmountOut, 10)
	amountInMaximum, _ := new(big.Int).SetString(s.AmountInMaximum, 10)
	fee, _ := new(big.Int).SetString(s.Fee, 10)
	sqrtPriceLimitX96, _ := new(big.Int).SetString(s.SqrtPriceLimitX96, 10)

	return SwapExactOutputSingleArgs{
		TokenIn:           s.TokenIn,
		AmountOut:         amountOut,
		TokenOut:          s.TokenOut,
		AmountInMaximum:   amountInMaximum,
		Fee:               fee,
		Recipient:         s.Recipient,
		SqrtPriceLimitX96: sqrtPriceLimitX96,
	}
}

func (s *SwapExactOutputSingleArgs) ConvertToJSONType() *JSONSwapExactOutputSingleArgs {
	return &JSONSwapExactOutputSingleArgs{
		TokenIn:           s.TokenIn,
		AmountOut:         s.AmountOut.String(),
		TokenOut:          s.TokenOut,
		AmountInMaximum:   s.AmountInMaximum.String(),
		Fee:               s.Fee.String(),
		Recipient:         s.Recipient,
		SqrtPriceLimitX96: s.SqrtPriceLimitX96.String(),
	}
}

func (s *SwapExactOutputSingleArgs) Decode(ctx context.Context, args map[string]interface{}) error {
	params, ok := args["params"].(struct {
		TokenIn           common.Address "json:\"tokenIn\""
		TokenOut          common.Address "json:\"tokenOut\""
		Fee               *big.Int       "json:\"fee\""
		Recipient         common.Address "json:\"recipient\""
		AmountOut         *big.Int       "json:\"amountOut\""
		AmountInMaximum   *big.Int       "json:\"amountInMaximum\""
		SqrtPriceLimitX96 *big.Int       "json:\"sqrtPriceLimitX96\""
	})
	if !ok {
		return fmt.Errorf("params is not of the expected type")
	}
	s.TokenIn = accounts.Address(params.TokenIn)
	s.TokenOut = accounts.Address(params.TokenOut)
	s.Fee = params.Fee
	s.Recipient = accounts.Address(params.Recipient)
	s.AmountOut = params.AmountOut
	s.AmountInMaximum = params.AmountInMaximum
	s.SqrtPriceLimitX96 = params.SqrtPriceLimitX96
	return nil
}

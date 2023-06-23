package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

const (
	swapExactInputSingle    = "swapExactInputSingle"
	swapExactOutputSingle   = "swapExactOutputSingle"
	swapExactInputMultihop  = "swapExactInputMultihop"
	swapExactOutputMultihop = "swapExactOutputMultihop"
	multicall               = "multicall"
)

func (u *UniswapClient) ProcessUniswapV3RouterTxs(ctx context.Context, tx MevTx) {
	switch tx.MethodName {
	case swapExactTokensForTokens:
		inputs := &SwapExactInputSingleArgs{}
		err := inputs.DecodeSwapExactInputSingle(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode swap exact input single args")
			return
		}
		// convert, get pricing data, run bin search
	case multicall:
	case swapExactInputSingle:
	case swapExactOutputSingle:
	case swapExactInputMultihop:
	case swapExactOutputMultihop:
	}
}

type SwapExactInputSingleArgs struct {
	TokenIn           accounts.Address `json:"tokenIn"`
	AmountIn          *big.Int         `json:"amountIn"`
	TokenOut          accounts.Address `json:"tokenOut"`
	AmountOutMinimum  *big.Int         `json:"amountOutMinimum"`
	Fee               *big.Int         `json:"fee"`
	Recipient         accounts.Address `json:"recipient"`
	SqrtPriceLimitX96 *big.Int         `json:"sqrtPriceLimitX96"`
}

func (s *SwapExactInputSingleArgs) DecodeSwapExactInputSingle(ctx context.Context, args map[string]interface{}) error {
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

/*
func (p *Permit2TransferFromParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[Permit2TransferFrom].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	token, err := ConvertToAddress(args["token"])
	if err != nil {
		return err
	}
	recipient, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	amount, err := ParseBigInt(args["amount"])
	if err != nil {
		return err
	}
	p.Token = token
	p.Recipient = recipient
	p.Amount = amount
	return nil
}

*/

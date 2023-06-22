package web3_client

import (
	"context"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

const (
	QuoterV1Address = "0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6"
	QuoterV2Address = "0x61fFE014bA17989E743c5F6cB21bF9697530B21e"

	quoteExactInput       = "quoteExactInput"
	quoteExactInputSingle = "quoteExactInputSingle"
)

type QuoteExactInputParams struct {
	TokenFeePath
	AmountIn          *big.Int
	SqrtPriceLimitX96 *big.Int
}

type QuoteExactInputSingleParams struct {
	TokenIn           accounts.Address `abi:"tokenIn"`
	TokenOut          accounts.Address `abi:"tokenOut"`
	Fee               *big.Int         `abi:"fee"`
	AmountIn          *big.Int         `abi:"amountIn"`
	SqrtPriceLimitX96 *big.Int         `abi:"sqrtPriceLimitX96"`
}

type UniswapAmountOutV3 struct {
	AmountOut               *big.Int
	SqrtPriceX96After       *big.Int
	InitializedTicksCrossed uint32
	GasEstimate             *big.Int
}

//func (u *UniswapClient) GetPoolV3QuoteFromQuoterV2(ctx context.Context, qp QuoteExactInputParams) (UniswapAmountOutV3, error) {
//	scInfo := &web3_actions.SendContractTxPayload{
//		SmartContractAddr: QuoterV2Address,
//		SendEtherPayload:  web3_actions.SendEtherPayload{},
//		ContractABI:       MustLoadQuoterV2Abi(),
//		MethodName:        quoteExactInput,
//		Params:            []interface{}{qp.TokenFeePath.Encode(), qp.AmountIn},
//	}
//
//	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
//	if err != nil {
//		return UniswapAmountOutV3{}, err
//	}
//	fmt.Println(resp)
//	return UniswapAmountOutV3{}, nil
//}

func (u *UniswapClient) GetPoolV3ExactInputSingleQuoteFromQuoterV2(ctx context.Context, qp QuoteExactInputSingleParams) (UniswapAmountOutV3, error) {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: QuoterV2Address,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       MustLoadQuoterV2Abi(),
		MethodName:        quoteExactInputSingle,
		Params:            []interface{}{qp},
	}
	qa := UniswapAmountOutV3{}
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return qa, err
	}
	for i, val := range resp {
		switch i {
		case 0:
			qa.AmountOut = val.(*big.Int)
		case 1:
			qa.SqrtPriceX96After = val.(*big.Int)
		case 2:
			qa.InitializedTicksCrossed = val.(uint32)
		case 3:
			qa.GasEstimate = val.(*big.Int)
		}
	}
	return qa, nil
}

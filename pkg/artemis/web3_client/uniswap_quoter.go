package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
)

const (
	QuoterV1Address = "0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6"
	QuoterV2Address = "0x61fFE014bA17989E743c5F6cB21bF9697530B21e"

	quoteExactInput = "quoteExactInput"
)

type QuoteExactInputSingleParams struct {
	TokenIn           accounts.Address
	TokenOut          accounts.Address
	Fee               uint64
	AmountIn          *big.Int
	SqrtPriceLimitX96 *big.Int
}

func (u *UniswapClient) GetPoolV3QuoteFromQuoterV2(ctx context.Context, poolV3 entities.Pool) (UniswapV2Pair, error) {
	qp := QuoteExactInputSingleParams{
		TokenIn:           poolV3.Token0.Address,
		TokenOut:          poolV3.Token1.Address,
		Fee:               uint64(poolV3.Fee),
		AmountIn:          big.NewInt(100),
		SqrtPriceLimitX96: big.NewInt(0),
	}
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: QuoterV2Address,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.PoolV3Abi,
		MethodName:        quoteExactInput,
		Params:            []interface{}{qp},
	}

	resp, err := u.Web3Client.CallFunctionWithArgs(ctx, scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	fmt.Println(resp)
	return UniswapV2Pair{}, nil
}

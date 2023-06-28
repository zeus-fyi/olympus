package web3_client

import (
	"context"
	"errors"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (u *UniswapClient) SingleReadMethodBigInt(ctx context.Context, methodName string, scInfo *web3_actions.SendContractTxPayload) (*big.Int, error) {
	scInfo.MethodName = methodName
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return &big.Int{}, err
	}
	if len(resp) == 0 {
		return &big.Int{}, errors.New("empty response")
	}
	bi, err := ParseBigInt(resp[0])
	if err != nil {
		return &big.Int{}, err
	}
	return bi, nil
}

func (u *UniswapClient) SingleReadMethodAddr(ctx context.Context, methodName string, scInfo *web3_actions.SendContractTxPayload) (accounts.Address, error) {
	scInfo.MethodName = methodName
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return accounts.Address{}, err
	}
	if len(resp) == 0 {
		return accounts.Address{}, errors.New("empty response")
	}
	addr, err := ConvertToAddress(resp[0])
	if err != nil {
		return accounts.Address{}, err
	}
	return addr, nil
}

func (p *UniswapV2Pair) sortTokens(tkn0, tkn1 accounts.Address) {
	token0Rep := big.NewInt(0).SetBytes(tkn0.Bytes())
	token1Rep := big.NewInt(0).SetBytes(tkn1.Bytes())

	if token0Rep.Cmp(token1Rep) > 0 {
		tkn0, tkn1 = tkn1, tkn0
	}
	p.Token0 = tkn0
	p.Token1 = tkn1
}

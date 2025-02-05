package web3_client

import (
	"context"
	"errors"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
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

func (u *UniswapClient) V2PairToPrices(ctx context.Context, bn uint64, pairAddr []accounts.Address) (*uniswap_pricing.UniswapV2Pair, error) {
	return uniswap_pricing.V2PairToPrices(ctx, bn, u.Web3Client.Web3Actions, pairAddr)
}

func (u *UniswapClient) GetPairContractFromFactory(ctx context.Context, addressOne, addressTwo string) accounts.Address {
	addrOne, addrTwo := StringsToAddresses(addressOne, addressTwo)
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: UniswapV2FactoryAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.FactoryAbi,
		MethodName:        "getPair",
		Params:            []interface{}{addrOne, addrTwo},
	}
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return accounts.Address{}
	}
	if len(resp) == 0 {
		return accounts.Address{}
	}
	pairAddr, err := ConvertToAddress(resp[0])
	if err != nil {
		return accounts.Address{}
	}
	return pairAddr
}

package web3_client

import (
	"context"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

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

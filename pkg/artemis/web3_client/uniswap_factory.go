package web3_client

import (
	"context"

	"github.com/gochain/gochain/v4/common"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

func (u *UniswapV2Client) GetPairContractFromFactory(ctx context.Context, addressOne, addressTwo string) common.Address {
	addrOne := common.HexToAddress(addressOne)
	addrTwo := common.HexToAddress(addressTwo)
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: UniswapV2FactoryAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.FactoryAbi,
		MethodName:        "getPair",
		Params:            []interface{}{addrOne, addrTwo},
	}
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return common.Address{}
	}
	if len(resp) == 0 {
		return common.Address{}
	}
	pairAddr, err := ConvertToAddress(resp[0])
	if err != nil {
		return common.Address{}
	}
	return pairAddr
}

package web3_client

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

const (
	Owner = "owner"
)

func (w *Web3Client) GetOwner(ctx context.Context, abiFile *abi.ABI, contractAddress string) (common.Address, error) {
	w.Dial()
	defer w.C.Close()
	payload := web3_actions.SendContractTxPayload{
		SmartContractAddr: contractAddress,
		ContractABI:       abiFile,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		MethodName:        Owner,
	}
	payload.Params = []interface{}{}
	owner, err := w.GetContractConst(ctx, &payload)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("GetOwner")
		return common.Address{}, err
	}
	return owner[0].(common.Address), err
}

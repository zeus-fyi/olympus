package web3_client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

const (
	execute = "execute"
)

type UniversalRouterExecCmd struct {
	Commands []UniversalRouterExecSubCmd `json:"commands"`
	Deadline *big.Int                    `json:"deadline"`
}

type UniversalRouterExecSubCmd struct {
	Command       string `json:"command"`
	CanRevert     bool   `json:"canRevert"`
	Inputs        []byte `json:"inputs"`
	DecodedInputs any    `json:"decodedInputs"`
}

func GetUniswapUniversalRouterAbiPayload(payload UniversalRouterExecCmd) web3_actions.SendContractTxPayload {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: UniswapUniversalRouterAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       UniversalRouterDecoder,
		MethodName:        execute,
		Params:            []interface{}{payload.Commands, payload.Deadline},
	}
	return params
}

func (u *UniswapClient) ExecUniswapUniversalRouterCmd(payload UniversalRouterExecCmd) (*types.Transaction, error) {
	scInfo := GetUniswapUniversalRouterAbiPayload(payload)
	// TODO implement better gas estimation
	scInfo.GasLimit = 3000000
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return nil, err
	}
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

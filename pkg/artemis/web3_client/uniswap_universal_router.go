package web3_client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

const (
	execute0 = "execute0"
	execute  = "execute"
)

type UniversalRouterExecCmd struct {
	Commands []UniversalRouterExecSubCmd    `json:"commands"`
	Deadline *big.Int                       `json:"deadline"`
	Payable  *web3_actions.SendEtherPayload `json:"payable,omitempty"`
}

type UniversalRouterExecSubCmd struct {
	Command       string `json:"command"`
	CanRevert     bool   `json:"canRevert"`
	Inputs        []byte `json:"inputs"`
	DecodedInputs any    `json:"decodedInputs"`
}

func GetUniswapUniversalRouterAbiPayload(payload *UniversalRouterExecParams) web3_actions.SendContractTxPayload {
	if payload == nil {
		payload = &UniversalRouterExecParams{}
		log.Warn().Msg("GetUniswapUniversalRouterAbiPayload: payload is nil")
		return web3_actions.SendContractTxPayload{}
	}
	payable := payload.Payable
	if payable == nil {
		payable = &web3_actions.SendEtherPayload{
			TransferArgs:   web3_actions.TransferArgs{},
			GasPriceLimits: web3_actions.GasPriceLimits{},
		}
	}
	fnParams := []interface{}{payload.Commands, payload.Inputs}
	methodName := execute
	if payload.Deadline != nil {
		methodName = execute0
		fnParams = append(fnParams, payload.Deadline.String())
	}
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: UniswapUniversalRouterAddress,
		SendEtherPayload:  *payable,
		ContractABI:       MustLoadUniversalRouterAbi(),
		MethodName:        methodName,
		Params:            fnParams,
	}
	return params
}

func (u *UniswapClient) ExecUniswapUniversalRouterCmd(payload UniversalRouterExecCmd) (*types.Transaction, error) {
	data, err := payload.EncodeCommands(ctx)
	if err != nil {
		log.Err(err).Msg("ExecUniswapUniversalRouterCmd: failed to encode commands")
		return nil, err
	}
	scInfo := GetUniswapUniversalRouterAbiPayload(data)
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

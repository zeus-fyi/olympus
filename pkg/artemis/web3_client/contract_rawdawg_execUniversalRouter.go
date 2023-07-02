package web3_client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

const (
	executeUniversalRouter = "executeUniversalRouter"
)

type UniversalRouterExecParams struct {
	Commands []byte                         `json:"commands"`
	Inputs   [][]byte                       `json:"inputs"`
	Deadline *big.Int                       `json:"deadline"`
	Payable  *web3_actions.SendEtherPayload `json:"payable,omitempty"`
}

func GetRawdawgUniversalRouterPayload(payload *UniversalRouterExecParams) web3_actions.SendContractTxPayload {
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
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: RawDawgAddr,
		SendEtherPayload:  *payable,
		ContractABI:       RawdawgAbi,
		MethodName:        executeUniversalRouter,
		Params:            []interface{}{payload.Commands, payload.Inputs, payload.Deadline},
	}
	return params
}

func (u *UniswapClient) ExecRawdawgUniversalRouterCmd(payload UniversalRouterExecCmd, to *artemis_trading_types.TradeOutcome) (*types.Transaction, error) {
	data, err := payload.EncodeCommands(ctx)
	if err != nil {
		log.Err(err).Msg("ExecRawdawgUniversalRouterCmd: failed to encode commands")
		return nil, err
	}
	scInfo := GetRawdawgUniversalRouterPayload(data)
	// TODO implement better gas estimation
	scInfo.GasLimit = 3000000
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return nil, err
	}
	if to != nil {
		to.AddTxHash(accounts.Hash(signedTx.Hash()))
	}
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		log.Err(err).Msg("ExecRawdawgUniversalRouterCmd: failed to send signed tx")
		return nil, err
	}
	return signedTx, nil
}

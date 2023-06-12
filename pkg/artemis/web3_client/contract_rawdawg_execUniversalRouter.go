package web3_client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
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

func GetRawdawgUniversalRouterPayload(tradeParams UniversalRouterExecParams) web3_actions.SendContractTxPayload {
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: UniswapUniversalRouterAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       RawdawgAbi,
		MethodName:        executeUniversalRouter,
		Params:            []interface{}{tradeParams.Commands, tradeParams.Inputs, tradeParams.Deadline},
	}
	return params
}

func (u *UniswapClient) ExecRawdawgUniversalRouterTrade(tradeParams UniversalRouterExecParams, to *TradeOutcome) (*types.Transaction, error) {
	scInfo := GetRawdawgUniversalRouterPayload(tradeParams)
	// TODO implement better gas estimation
	scInfo.GasLimit = 3000000
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return nil, err
	}
	to.AddTxHash(accounts.Hash(signedTx.Hash()))
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

package web3_client

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

const (
	Aggregate3      = "aggregate3"
	Aggregate3Value = "aggregate3Value"
)

type Multicall3 struct {
	Name          string         `json:"name"`
	Target        common.Address `json:"target"`
	AbiFile       *abi.ABI       `json:"abiFile"`
	CallData      []byte         `json:"callData"`
	DecodedInputs []any          `json:"decodedInputs,omitempty"`
}

// Multicall3Result is an auto generated low-level Go binding around an user-defined struct.
type Multicall3Result struct {
	Success    bool   `json:"success"`
	ReturnData []byte `json:"returnData"`
}

func CreateMulticall3Payload(ctx context.Context, calls []Multicall3) (web3_actions.SendContractTxPayload, error) {
	payload := web3_actions.SendContractTxPayload{
		SmartContractAddr: artemis_trading_constants.Multicall3Address,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       artemis_oly_contract_abis.MultiCall3,
		MethodName:        Aggregate3,
		Params:            []interface{}{},
	}
	for _, ca := range calls {
		inputs, err := ca.AbiFile.Methods[ca.Name].Inputs.Pack(ca.DecodedInputs...)
		if err != nil {
			return web3_actions.SendContractTxPayload{}, err
		}
		ca.CallData = inputs
		payload.Params = append(payload.Params, ca.Target, ca.CallData)
	}
	return payload, nil
}

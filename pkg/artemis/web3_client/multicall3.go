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

type MultiCallElement struct {
	Name          string        `json:"name"`
	AbiFile       *abi.ABI      `json:"abiFile"`
	Call          Call          `json:"callData"`
	DecodedInputs []interface{} `json:"decodedInputs,omitempty"`
}
type Call struct {
	Target       common.Address `abi:"target"`
	AllowFailure bool           `abi:"allowFailure"`
	Data         []byte         `abi:"callData"`
}

type Multicall3 struct {
	Calls []Call `json:"calls"`
}

// Multicall3Result is an auto generated low-level Go binding around an user-defined struct.
type Multicall3Result struct {
	Success    bool   `json:"success"`
	ReturnData []byte `json:"returnData"`
}

func CreateMulticall3Payload(ctx context.Context, calls []MultiCallElement) (web3_actions.SendContractTxPayload, error) {
	callInputs := make([]Call, len(calls))
	for i, _ := range calls {
		method := calls[i].AbiFile.Methods[calls[i].Name]
		targetName := method.ID
		inputs, err := method.Inputs.Pack(calls[i].DecodedInputs...)
		if err != nil {
			return web3_actions.SendContractTxPayload{}, err
		}
		inputs = append(targetName[:], inputs[:]...)
		callInputs[i] = Call{
			Target:       calls[i].Call.Target,
			AllowFailure: calls[i].Call.AllowFailure,
			Data:         inputs,
		}
	}
	return makeAggregate3Payload(callInputs), nil
}

func makeAggregate3Payload(calls []Call) web3_actions.SendContractTxPayload {
	payload := web3_actions.SendContractTxPayload{
		SmartContractAddr: artemis_trading_constants.Multicall3Address,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       artemis_oly_contract_abis.MultiCall3,
		MethodName:        Aggregate3,
		Params:            []interface{}{calls},
	}
	return payload
}

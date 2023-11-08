package artemis_multicall

import (
	"context"
	"strings"

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
	Calls   []MultiCallElement `json:"calls"`
	Results []Multicall3Result `json:"results"`
}

// Multicall3Result is an auto generated low-level Go binding around an user-defined struct.
type Multicall3Result struct {
	Success           bool          `json:"success"`
	ReturnData        []byte        `json:"returnData"`
	DecodedReturnData []interface{} `json:"decodedReturnData,omitempty"`
}

func (m *Multicall3) PackAndCall(ctx context.Context, wc web3_actions.Web3Actions) ([]Multicall3Result, error) {
	payload, err := CreateMulticall3Payload(ctx, m.Calls)
	if err != nil {
		return nil, err
	}
	resp, err := wc.CallConstantFunction(ctx, &payload)
	if err != nil {
		return nil, err
	}
	encDataResponses := resp[0].([]struct {
		Success    bool    "json:\"success\""
		ReturnData []uint8 "json:\"returnData\""
	})
	results := make([]Multicall3Result, len(encDataResponses))
	for i, _ := range encDataResponses {
		encData := encDataResponses[i]
		results[i] = Multicall3Result{
			Success:    encData.Success,
			ReturnData: encData.ReturnData,
		}
		if encData.Success {
			decoded, derr := m.Calls[i].AbiFile.Methods[m.Calls[i].Name].Outputs.UnpackValues(encData.ReturnData)
			if derr != nil {
				return nil, derr
			}
			results[i].DecodedReturnData = decoded
		}
	}
	m.Results = results
	return results, err
}

func UnpackMultiCall(ctx context.Context, resp []interface{}, calls []MultiCallElement) ([]Multicall3Result, error) {
	encDataResponses := resp[0].([]struct {
		Success    bool    "json:\"success\""
		ReturnData []uint8 "json:\"returnData\""
	})
	results := make([]Multicall3Result, len(encDataResponses))
	for i, _ := range encDataResponses {
		encData := encDataResponses[i]
		results[i] = Multicall3Result{
			Success:    encData.Success,
			ReturnData: encData.ReturnData,
		}
		if encData.Success {
			decoded, derr := calls[i].AbiFile.Methods[calls[i].Name].Outputs.UnpackValues(encData.ReturnData)
			if derr != nil {
				return nil, derr
			}
			results[i].DecodedReturnData = decoded
		}
	}
	return results, nil
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

func UrCreateMulticall3Payload(ctx context.Context, calls []MultiCallElement, ur web3_actions.SendContractTxPayload) (web3_actions.SendContractTxPayload, error) {
	callInputs := make([]Call, len(calls))
	for i, _ := range calls {
		if strings.HasPrefix(calls[i].Name, "execute") {
			tmps := calls[i].Call.Data
			callInputs[i] = Call{
				Target:       calls[i].Call.Target,
				AllowFailure: calls[i].Call.AllowFailure,
				Data:         tmps,
			}
		} else {
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

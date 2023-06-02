package web3_client

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
)

func LoadSwapAbi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.SwapABI))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func LoadSwapAbiPayload(pairContractAddr string) (web3_actions.SendContractTxPayload, string, error) {
	abf := LoadSwapAbi()
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: pairContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractFile:      "",
		ContractABI:       abf,
		MethodName:        swap,
		Params:            []interface{}{},
	}
	return params, "", nil
}

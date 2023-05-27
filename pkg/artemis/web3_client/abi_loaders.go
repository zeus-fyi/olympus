package web3_client

import (
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func LoadERC20Abi() *abi.ABI {
	readAbi, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.ERC20ABI))
	if err != nil {
		panic(err)
	}
	return readAbi
}

func LoadERC20AbiPayload() (web3_actions.SendContractTxPayload, string, error) {
	fp := filepaths.Path{
		PackageName: "",
		DirIn:       "./contract_abis",
		FnIn:        "Token.json",
	}
	fi := fp.ReadFileInPath()
	m := map[string]interface{}{}
	err := json.Unmarshal(fi, &m)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	abiInput := m["abi"]
	b, err := json.Marshal(abiInput)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	abf := &abi.ABI{}
	err = abf.UnmarshalJSON(b)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, "", err
	}
	params := web3_actions.SendContractTxPayload{
		SendEtherPayload: web3_actions.SendEtherPayload{},
		ContractABI:      abf,
		Params:           []interface{}{},
	}
	return params, m["bytecode"].(string), nil
}

func LoadERC20DeployedByteCode() (string, error) {
	fp := filepaths.Path{
		PackageName: "",
		DirIn:       "./contract_abis",
		FnIn:        "Token.json",
	}
	fi := fp.ReadFileInPath()
	m := map[string]interface{}{}
	err := json.Unmarshal(fi, &m)
	if err != nil {
		return "", err
	}
	return m["deployedBytecode"].(string), nil
}

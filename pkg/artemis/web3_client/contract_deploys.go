package web3_client

import (
	"context"
	"encoding/json"

	"github.com/gochain/gochain/v4/accounts/abi"
	web3_types "github.com/zeus-fyi/gochain/web3/types"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

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

func (w *Web3Client) DeployERC20Token(ctx context.Context, bc string, scParams web3_actions.SendContractTxPayload) (*web3_types.Transaction, error) {
	w.Dial()
	defer w.Close()

	tx, err := w.DeployContract(ctx, bc, scParams)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

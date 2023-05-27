package web3_client

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

func LoadUniswapFactoryAbiPayload() (web3_actions.SendContractTxPayload, string, error) {
	fp := filepaths.Path{
		PackageName: "",
		DirIn:       "./contract_abis",
		FnIn:        "UniswapV2Factory.json",
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

func (w *Web3Client) DeployERC20Token(ctx context.Context, bc string, scParams web3_actions.SendContractTxPayload) (*types.Transaction, error) {
	w.Dial()
	defer w.Close()
	tx, err := w.DeployContract(ctx, bc, scParams)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

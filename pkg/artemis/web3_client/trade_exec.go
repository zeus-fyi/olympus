package web3_client

import (
	"encoding/json"
	"math/big"

	"github.com/gochain/gochain/v4/accounts/abi"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

const SwapContractAddr = ""

func LoadSwapAbiPayload() (web3_actions.SendContractTxPayload, string, error) {
	fp := filepaths.Path{
		PackageName: "",
		DirIn:       "./contract_abis",
		FnIn:        "RawSwap.json",
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
		SmartContractAddr: SwapContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractFile:      "",
		ContractABI:       abf,
		MethodName:        "executeSwap",
		Params:            []interface{}{},
	}
	return params, m["bytecode"].(string), nil
}

func (w *Web3Client) ExecSwapTrade(tf TradeExecutionFlowInBigInt) (*big.Int, *big.Int) {
	sellAmount := big.NewInt(0)
	maxProfit := big.NewInt(0)
	paramsTx, _, err := LoadSwapAbiPayload()
	if err != nil {
		return sellAmount, maxProfit
	}
	// Pair address in contract
	pairContract := ""
	paramsTx.Params = []interface{}{pairContract, tf.FrontRunTrade.AmountIn, tf.FrontRunTrade.AmountOut}
	return sellAmount, maxProfit
}

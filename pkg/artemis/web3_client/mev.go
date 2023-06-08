package web3_client

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

type MevSmartContractTxMap struct {
	SmartContractAddr string
	Abi               *abi.ABI
	Txs               []MevTx
	Filter            *strings_filter.FilterOpts
}

type MevTx struct {
	UserAddr    string
	MethodName  string
	Args        map[string]interface{}
	Order       string
	TxPoolQueue map[string]*types.Transaction
	Tx          *types.Transaction
}

package web3_client

import (
	"context"
	"fmt"
	"strings"

	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

const (
	UniswapV2RouterAddress = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	swapExactETHForTokens  = "swapExactETHForTokens"
)

// TODO https://docs.uniswap.org/contracts/v2/reference/smart-contracts/router-02

type UniswapV2Client struct {
	MevSmartContractTxMap
}

func InitUniswapV2Client(ctx context.Context) UniswapV2Client {
	abiFile, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniswapV2RouterABI))
	if err != nil {
		panic(err)
	}
	f := strings_filter.FilterOpts{
		DoesNotStartWithThese: nil,
		StartsWithThese:       []string{},
		DoesNotInclude:        nil,
	}
	return UniswapV2Client{MevSmartContractTxMap{
		SmartContractAddr: UniswapV2RouterAddress,
		Abi:               abiFile,
		MethodTxMap:       map[string]MevTx{},
		Txs:               []MevTx{},
		Filter:            &f,
	}}
}

// TODO finish this
func (u *UniswapV2Client) GetAllTradeMethods() []string {
	return []string{
		swapExactETHForTokens,
	}
}

func (u *UniswapV2Client) ProcessTxs() {
	for methodName, tx := range u.MethodTxMap {
		fmt.Println(tx)
		switch methodName {
		case swapExactETHForTokens:
			u.SwapExactETHForTokens(tx.Args)
		}
	}
}

func (u *UniswapV2Client) SwapExactETHForTokens(args map[string]interface{}) {
	fmt.Println("SwapExactETHForTokens", args)
}

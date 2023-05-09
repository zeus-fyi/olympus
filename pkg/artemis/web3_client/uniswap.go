package web3_client

import (
	"context"
	"fmt"
	"strings"

	"github.com/gochain/gochain/v4/common"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

const (
	UniswapV2RouterAddress       = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	addLiquidity                 = "addLiquidity"
	addLiquidityETH              = "addLiquidityETH"
	removeLiquidity              = "removeLiquidity"
	removeLiquidityETH           = "removeLiquidityETH"
	removeLiquidityWithPermit    = "removeLiquidityWithPermit"
	removeLiquidityETHWithPermit = "removeLiquidityETHWithPermit"
	swapExactTokensForTokens     = "swapExactTokensForTokens"
	swapTokensForExactTokens     = "swapTokensForExactTokens"
	swapExactETHForTokens        = "swapExactETHForTokens"
	swapTokensForExactETH        = "swapTokensForExactETH"
	swapExactTokensForETH        = "swapExactTokensForETH"
	swapETHForExactTokens        = "swapETHForExactTokens"
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

func (u *UniswapV2Client) GetAllTradeMethods() []string {
	return []string{
		swapExactETHForTokens,
		addLiquidity,
		addLiquidityETH,
		removeLiquidity,
		removeLiquidityETH,
		removeLiquidityWithPermit,
		removeLiquidityETHWithPermit,
		swapExactTokensForTokens,
		swapTokensForExactTokens,
		swapExactETHForTokens,
		swapTokensForExactETH,
		swapExactTokensForETH,
		swapETHForExactTokens,
	}
}

func (u *UniswapV2Client) ProcessTxs() {
	for methodName, tx := range u.MethodTxMap {
		fmt.Println(tx)
		switch methodName {
		case addLiquidity:
			//u.AddLiquidity(tx.Args)
		case addLiquidityETH:
			//u.AddLiquidityETH(tx.Args)
		case removeLiquidity:
			//u.RemoveLiquidity(tx.Args)
		case removeLiquidityETH:
			//u.RemoveLiquidityETH(tx.Args)
		case removeLiquidityWithPermit:
			//u.RemoveLiquidityWithPermit(tx.Args)
		case removeLiquidityETHWithPermit:
			//u.RemoveLiquidityETHWithPermit(tx.Args)
		case swapExactTokensForTokens:
			u.SwapExactTokensForTokens(tx.Args)
		case swapTokensForExactTokens:
			u.SwapTokensForExactTokens(tx.Args)
		case swapExactETHForTokens:
			u.SwapExactETHForTokens(tx.Args)
		case swapTokensForExactETH:
			u.SwapTokensForExactETH(tx.Args)
		case swapExactTokensForETH:
			u.SwapExactTokensForETH(tx.Args)
		case swapETHForExactTokens:
			u.SwapETHForExactTokens(tx.Args)
		}
	}
}

func (u *UniswapV2Client) SwapExactTokensForTokens(args map[string]interface{}) {
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		return
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		return
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return
	}
	to := common.HexToAddress(args["to"].(string))
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		return
	}
	fmt.Println("amountIn:", amountIn)
	fmt.Println("amountOutMin:", amountOutMin)
	fmt.Println("path:", path)
	fmt.Println("to:", to)
	fmt.Println("deadline:", deadline)
}

func (u *UniswapV2Client) SwapTokensForExactTokens(args map[string]interface{}) {
	amountOut, err := ParseBigInt(args["amountOut"])
	if err != nil {
		return
	}
	amountInMax, err := ParseBigInt(args["amountInMax"])
	if err != nil {
		return
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return
	}
	to := common.HexToAddress(args["to"].(string))
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		return
	}
	fmt.Println("amountOut:", amountOut)
	fmt.Println("amountInMax:", amountInMax)
	fmt.Println("path:", path)
	fmt.Println("to:", to)
	fmt.Println("deadline:", deadline)
}

func (u *UniswapV2Client) SwapExactETHForTokens(args map[string]interface{}) {
	fmt.Println("SwapExactETHForTokens", args)
}

func (u *UniswapV2Client) SwapTokensForExactETH(args map[string]interface{}) {
	fmt.Println("SwapTokensForExactETH", args)
}

func (u *UniswapV2Client) SwapExactTokensForETH(args map[string]interface{}) {
	fmt.Println("SwapExactTokensForETH", args)
}

func (u *UniswapV2Client) SwapETHForExactTokens(args map[string]interface{}) {
	fmt.Println("SwapExactTokensForETH", args)
}

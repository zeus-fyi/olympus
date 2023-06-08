package web3_client

import (
	"context"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

const (
	UniswapUniversalRouterAddress = "0xEf1c6E67703c7BD7107eed8303Fbe6EC2554BF6B"
	UniswapV2FactoryAddress       = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
	UniswapV2RouterAddress        = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"

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
	getAmountsOutFrontRunTrade   = "getAmountsOutFrontRunTrade"
	getAmountsOut                = "getAmountsOut"
	getAmountsIn                 = "getAmountsIn"

	execSmartContractTradingSwap = "executeSwap"
	swap                         = "swap"
	swapFrontRun                 = "swapFrontRun"
	swapSandwich                 = "swapSandwich"
)

type UniswapClient struct {
	mu                                   sync.Mutex
	Web3Client                           Web3Client
	UniversalRouterSmartContractAddr     string
	FactorySmartContractAddr             string
	RouterSmartContractAddr              string
	PairAbi                              *abi.ABI
	ERC20Abi                             *abi.ABI
	FactoryAbi                           *abi.ABI
	UniversalRouterAbi                   *abi.ABI
	PrintDetails                         bool
	PrintOn                              bool
	PrintLocal                           bool
	DebugPrint                           bool
	MevSmartContractTxMapUniversalRouter MevSmartContractTxMap
	MevSmartContractTxMap
	*TradeAnalysisReport
	Path                                filepaths.Path
	BlockNumber                         *big.Int
	Trades                              []artemis_autogen_bases.EthMempoolMevTx
	chronos                             chronos.Chronos
	SwapExactTokensForTokensParamsSlice []SwapExactTokensForTokensParams
	SwapTokensForExactTokensParamsSlice []SwapTokensForExactTokensParams
	SwapExactETHForTokensParamsSlice    []SwapExactETHForTokensParams
	SwapTokensForExactETHParamsSlice    []SwapTokensForExactETHParams
	SwapExactTokensForETHParamsSlice    []SwapExactTokensForETHParams
	SwapETHForExactTokensParamsSlice    []SwapETHForExactTokensParams
}

func InitUniswapClient(ctx context.Context, w Web3Client) UniswapClient {
	erc20AbiFile, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.ERC20ABI))
	if err != nil {
		panic(err)
	}
	factoryAbiFile, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniswapV2FactoryABI))
	if err != nil {
		panic(err)
	}
	pairAbiFile, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniswapV2PairAbi))
	if err != nil {
		panic(err)
	}
	f := strings_filter.FilterOpts{
		DoesNotStartWithThese: nil,
		StartsWithThese:       []string{"swap"},
		Contains:              "",
		DoesNotInclude:        []string{"supportingFeeOnTransferTokens"},
	}
	return UniswapClient{
		Web3Client:                       w,
		chronos:                          chronos.Chronos{},
		FactorySmartContractAddr:         UniswapV2FactoryAddress,
		RouterSmartContractAddr:          UniswapV2RouterAddress,
		UniversalRouterSmartContractAddr: UniswapUniversalRouterAddress,
		FactoryAbi:                       factoryAbiFile,
		ERC20Abi:                         erc20AbiFile,
		PairAbi:                          pairAbiFile,
		UniversalRouterAbi:               MustLoadUniversalRouterAbi(),
		MevSmartContractTxMapUniversalRouter: MevSmartContractTxMap{
			SmartContractAddr: UniswapUniversalRouterAddress,
			Abi:               MustLoadUniversalRouterAbi(),
			Txs:               []MevTx{},
		},
		MevSmartContractTxMap: MevSmartContractTxMap{
			SmartContractAddr: UniswapV2RouterAddress,
			Abi:               MustLoadUniswapV2RouterABI(),
			Txs:               []MevTx{},
			Filter:            &f,
		},
		TradeAnalysisReport:                 &TradeAnalysisReport{},
		SwapExactTokensForTokensParamsSlice: []SwapExactTokensForTokensParams{},
		SwapTokensForExactTokensParamsSlice: []SwapTokensForExactTokensParams{},
		SwapExactETHForTokensParamsSlice:    []SwapExactETHForTokensParams{},
		SwapTokensForExactETHParamsSlice:    []SwapTokensForExactETHParams{},
		SwapExactTokensForETHParamsSlice:    []SwapExactTokensForETHParams{},
		SwapETHForExactTokensParamsSlice:    []SwapETHForExactTokensParams{},
	}
}

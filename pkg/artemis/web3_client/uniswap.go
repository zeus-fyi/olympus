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
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/web3/signing_automation/ethereum"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

const (
	UniswapUniversalRouterAddressNew = "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	UniswapUniversalRouterAddressOld = "0xEf1c6E67703c7BD7107eed8303Fbe6EC2554BF6B"
	UniswapV2FactoryAddress          = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
	UniswapV2Router02Address         = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	UniswapV2Router01Address         = "0xf164fC0Ec4E93095b804a4795bBe1e041497b92a"

	UniswapV3Router01Address = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	UniswapV3Router02Address = "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45"

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
	mu                                      sync.Mutex
	Web3Client                              Web3Client
	UniversalRouterSmartContractAddr        string
	FactorySmartContractAddr                string
	RouterSmartContractAddr                 string
	PoolV3Abi                               *abi.ABI
	PairAbi                                 *abi.ABI
	ERC20Abi                                *abi.ABI
	FactoryAbi                              *abi.ABI
	UniversalRouterAbi                      *abi.ABI
	PrintDetails                            bool
	PrintOn                                 bool
	PrintLocal                              bool
	DebugPrint                              bool
	TestMode                                bool
	MevSmartContractTxMapV3SwapRouterV1     MevSmartContractTxMap
	MevSmartContractTxMapV3SwapRouterV2     MevSmartContractTxMap
	MevSmartContractTxMapUniversalRouterNew MevSmartContractTxMap
	MevSmartContractTxMapUniversalRouterOld MevSmartContractTxMap
	MevSmartContractTxMapV2Router01         MevSmartContractTxMap
	MevSmartContractTxMapV2Router02         MevSmartContractTxMap
	*TradeAnalysisReport
	Path                                                          filepaths.Path
	BlockNumber                                                   *big.Int
	Trades                                                        []artemis_autogen_bases.EthMempoolMevTx
	chronos                                                       chronos.Chronos
	SwapExactTokensForTokensParamsSlice                           []SwapExactTokensForTokensParams
	SwapTokensForExactTokensParamsSlice                           []SwapTokensForExactTokensParams
	SwapExactETHForTokensParamsSlice                              []SwapExactETHForTokensParams
	SwapTokensForExactETHParamsSlice                              []SwapTokensForExactETHParams
	SwapExactTokensForETHParamsSlice                              []SwapExactTokensForETHParams
	SwapETHForExactTokensParamsSlice                              []SwapETHForExactTokensParams
	SwapExactTokensForETHSupportingFeeOnTransferTokensParamsSlice []SwapExactTokensForETHSupportingFeeOnTransferTokensParams
	SwapExactETHForTokensSupportingFeeOnTransferTokensParamsSlice []SwapExactETHForTokensSupportingFeeOnTransferTokensParams
	SwapExactTokensForTokensSupportingFeeOnTransferTokensSlice    []SwapExactTokensForTokensSupportingFeeOnTransferTokensParams
}

func InitUniswapClient(ctx context.Context, w Web3Client) UniswapClient {
	erc20AbiFile, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.Erc20Abi))
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
		DoesNotInclude:        []string{""},
	}
	return UniswapClient{
		Web3Client:                       w,
		chronos:                          chronos.Chronos{},
		FactorySmartContractAddr:         UniswapV2FactoryAddress,
		RouterSmartContractAddr:          UniswapV2Router02Address,
		UniversalRouterSmartContractAddr: UniswapUniversalRouterAddressNew,
		FactoryAbi:                       factoryAbiFile,
		ERC20Abi:                         erc20AbiFile,
		PoolV3Abi:                        artemis_oly_contract_abis.MustLoadPoolV3Abi(),
		PairAbi:                          pairAbiFile,
		UniversalRouterAbi:               artemis_oly_contract_abis.MustLoadNewUniversalRouterAbi(),
		MevSmartContractTxMapUniversalRouterNew: MevSmartContractTxMap{
			SmartContractAddr: UniswapUniversalRouterAddressNew,
			Abi:               artemis_oly_contract_abis.MustLoadNewUniversalRouterAbi(),
			Txs:               []MevTx{},
		},
		MevSmartContractTxMapUniversalRouterOld: MevSmartContractTxMap{
			SmartContractAddr: UniswapUniversalRouterAddressOld,
			Abi:               artemis_oly_contract_abis.MustLoadOldUniversalRouterAbi(),
			Txs:               []MevTx{},
		},
		MevSmartContractTxMapV2Router02: MevSmartContractTxMap{
			SmartContractAddr: UniswapV2Router02Address,
			Abi:               artemis_oly_contract_abis.MustLoadUniswapV2Router02ABI(),
			Txs:               []MevTx{},
			Filter:            &f,
		},
		MevSmartContractTxMapV2Router01: MevSmartContractTxMap{
			SmartContractAddr: UniswapV2Router01Address,
			Abi:               artemis_oly_contract_abis.MustLoadUniswapV2Router01ABI(),
			Txs:               []MevTx{},
			Filter:            &f,
		},
		MevSmartContractTxMapV3SwapRouterV1: MevSmartContractTxMap{
			SmartContractAddr: UniswapV3Router01Address,
			Abi:               artemis_oly_contract_abis.MustLoadUniswapV3Swap1RouterAbi(),
			Txs:               []MevTx{},
		},
		MevSmartContractTxMapV3SwapRouterV2: MevSmartContractTxMap{
			SmartContractAddr: UniswapV3Router02Address,
			Abi:               artemis_oly_contract_abis.MustLoadUniswapV3Swap2RouterAbi(),
			Txs:               []MevTx{},
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

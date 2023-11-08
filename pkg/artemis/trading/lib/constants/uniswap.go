package artemis_trading_constants

import (
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

const (
	UniswapUniversalRouterAddressNew = "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	UniswapUniversalRouterAddressOld = "0xEf1c6E67703c7BD7107eed8303Fbe6EC2554BF6B"
	UniswapV2Router02Address         = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	UniswapV2Router01Address         = "0xf164fC0Ec4E93095b804a4795bBe1e041497b92a"

	UniswapV3Router01Address = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	UniswapV3Router02Address = "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45"

	UniswapV3FactoryAddress = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	UniswapV2FactoryAddress = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"

	UniswapQuoterAddress = "0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6"
)

const (
	UniversalRouterSender    = "0x0000000000000000000000000000000000000001"
	UniversalRouterRecipient = "0x0000000000000000000000000000000000000002"
)

const (
	V2SwapExactIn                                         = "V2_SWAP_EXACT_IN"
	V2SwapExactOut                                        = "V2_SWAP_EXACT_OUT"
	V3SwapExactIn                                         = "V3_SWAP_EXACT_IN"
	V3SwapExactOut                                        = "V3_SWAP_EXACT_OUT"
	SwapExactTokensForETH                                 = "swapExactTokensForETH"
	SwapETHForExactTokens                                 = "swapETHForExactTokens"
	SwapTokensForExactETH                                 = "swapTokensForExactETH"
	SwapTokensForExactTokens                              = "swapTokensForExactTokens"
	SwapExactETHForTokens                                 = "swapExactETHForTokens"
	SwapExactTokensForTokens                              = "swapExactTokensForTokens"
	SwapExactTokensForETHSupportingFeeOnTransferTokens    = "swapExactTokensForETHSupportingFeeOnTransferTokens"
	SwapExactETHForTokensSupportingFeeOnTransferTokens    = "swapExactETHForTokensSupportingFeeOnTransferTokens"
	SwapExactTokensForTokensSupportingFeeOnTransferTokens = "swapExactTokensForTokensSupportingFeeOnTransferTokens"
	Permit2TransferFrom                                   = "PERMIT2_TRANSFER_FROM"
	Permit2PermitBatch                                    = "PERMIT2_PERMIT_BATCH"
	Permit2Permit                                         = "PERMIT2_PERMIT"
	Permit2TransferFromBatch                              = "PERMIT2_TRANSFER_FROM_BATCH"

	Multicall = "multicall"
	Execute0  = "execute0"
	Execute1  = "execute1"

	Execute = "execute"

	Sweep      = "SWEEP"
	PayPortion = "PAY_PORTION"
	Transfer   = "TRANSFER"
	UnwrapWETH = "UNWRAP_WETH"
	WrapETH    = "WRAP_ETH"
)

var (
	UniswapUniversalRouterNewAddressAccount = accounts.HexToAddress(UniswapUniversalRouterAddressNew)
	UniswapV3FactoryAddressAccount          = accounts.HexToAddress(UniswapV3FactoryAddress)
	UniswapV2FactoryAddressAccount          = accounts.HexToAddress(UniswapV2FactoryAddress)

	UniversalRouterSenderAddress   = accounts.HexToAddress("0x0000000000000000000000000000000000000001")
	UniversalRouterReceiverAddress = accounts.HexToAddress("0x0000000000000000000000000000000000000002")
)

var (
	UniversalRouterNewAbi = artemis_oly_contract_abis.MustLoadNewUniversalRouterAbi()
)

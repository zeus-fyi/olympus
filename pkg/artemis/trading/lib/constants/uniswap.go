package artemis_trading_constants

import (
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

const (
	UniswapUniversalRouterAddressNew = "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	UniswapUniversalRouterAddressOld = "0xEf1c6E67703c7BD7107eed8303Fbe6EC2554BF6B"

	UniswapV3FactoryAddress = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	UniswapV2FactoryAddress = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
)

const (
	UniversalRouterSender    = "0x0000000000000000000000000000000000000001"
	UniversalRouterRecipient = "0x0000000000000000000000000000000000000002"
)

const (
	V2SwapExactIn            = "V2_SWAP_EXACT_IN"
	V2SwapExactOut           = "V2_SWAP_EXACT_OUT"
	SwapExactTokensForTokens = "swapExactTokensForTokens"

	Permit2TransferFrom      = "PERMIT2_TRANSFER_FROM"
	Permit2PermitBatch       = "PERMIT2_PERMIT_BATCH"
	Permit2Permit            = "PERMIT2_PERMIT"
	Permit2TransferFromBatch = "PERMIT2_TRANSFER_FROM_BATCH"

	Execute0 = "execute0"
	Execute  = "execute"

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

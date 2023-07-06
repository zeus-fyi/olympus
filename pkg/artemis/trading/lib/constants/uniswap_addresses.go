package artemis_trading_constants

import "github.com/zeus-fyi/gochain/web3/accounts"

const (
	UniswapV3FactoryAddress = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	UniswapV2FactoryAddress = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
)

const (
	V2SwapExactIn  = "V2_SWAP_EXACT_IN"
	V2SwapExactOut = "V2_SWAP_EXACT_OUT"
)

var (
	UniswapV3FactoryAddressAccount = accounts.HexToAddress(UniswapV3FactoryAddress)
	UniswapV2FactoryAddressAccount = accounts.HexToAddress(UniswapV2FactoryAddress)
)

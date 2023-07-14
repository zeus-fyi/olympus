package artemis_pricing_utils

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

func ApplyTransferTax(tokenAddress accounts.Address, amount *big.Int) *big.Int {
	if artemis_trading_cache.TokenMap == nil {
		return amount
	}
	num := artemis_trading_cache.TokenMap[tokenAddress.String()].TransferTaxNumerator
	denom := artemis_trading_cache.TokenMap[tokenAddress.String()].TransferTaxDenominator
	if num == nil || denom == nil {
		panic("numerator or denominator is nil")
		return amount
	}
	if *num == 1 && *denom == 1 {
		return amount
	}
	return artemis_eth_units.ApplyTransferTax(amount, *num, *denom)
}

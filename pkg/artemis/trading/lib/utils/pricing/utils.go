package artemis_pricing_utils

import (
	"errors"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

func ApplyTransferTax(tokenAddress accounts.Address, amount *big.Int) (*big.Int, error) {
	if artemis_trading_cache.TokenMap == nil {
		return amount, errors.New("TokenMap is nil")
	}
	num := artemis_trading_cache.TokenMap[tokenAddress.String()].TransferTaxNumerator
	denom := artemis_trading_cache.TokenMap[tokenAddress.String()].TransferTaxDenominator
	if num == nil || denom == nil {
		return amount, errors.New("transfer tax is nil")
	}
	if *num == 1 && *denom == 1 {
		return amount, nil
	}
	return artemis_eth_units.ApplyTransferTax(amount, *num, *denom), nil
}

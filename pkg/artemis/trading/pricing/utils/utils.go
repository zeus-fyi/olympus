package artemis_pricing_utils

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
)

func ApplyTransferTax(tokenAddress accounts.Address, amount *big.Int) *big.Int {
	if artemis_trading_cache.TokenMap == nil {
		return amount
	}
	num := artemis_trading_cache.TokenMap[tokenAddress.String()].TransferTaxNumerator
	denom := artemis_trading_cache.TokenMap[tokenAddress.String()].TransferTaxDenominator
	if num == nil || denom == nil {
		return amount
	}
	if *num == 1 && *denom == 1 {
		return amount
	}
	transferTax := uniswap_core_entities.NewPercent(new(big.Int).SetInt64(int64(*num)), new(big.Int).SetInt64(int64(*denom)))
	transferFee := new(big.Int).Mul(amount, transferTax.Numerator)
	transferFee = transferFee.Div(transferFee, transferTax.Denominator)
	adjustedOut := new(big.Int).Sub(amount, transferFee)
	return adjustedOut
}

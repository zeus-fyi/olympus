package artemis_trading_auxiliary

import (
	"math/big"

	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

func (a *AuxiliaryTradingUtils) maxTradeSize() *big.Int {
	gweiInEther := artemis_eth_units.GweiPerEth
	return artemis_eth_units.GweiMultiple(gweiInEther / 4)
}

func (a *AuxiliaryTradingUtils) isProfitHigherThanGasFee() bool {
	return false
}

func (a *AuxiliaryTradingUtils) isTradingEnabledOnToken() bool {
	return false
}

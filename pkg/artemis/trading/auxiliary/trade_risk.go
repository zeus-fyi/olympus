package artemis_trading_auxiliary

import (
	"math/big"

	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
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

// in sandwich trade the tokenIn on the first trade is the profit currency
func (a *AuxiliaryTradingUtils) isProfitTokenAcceptable(tf *web3_client.TradeExecutionFlow) bool {
	wethAddr := a.getChainSpecificWETH()
	if tf.FrontRunTrade.AmountInAddr.String() != wethAddr.String() {
		return false
	}
	return true
}

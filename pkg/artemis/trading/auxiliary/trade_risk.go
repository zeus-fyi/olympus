package artemis_trading_auxiliary

import (
	"context"
	"errors"
	"math/big"

	"github.com/rs/zerolog/log"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) maxTradeSize() *big.Int {
	gweiInEther := artemis_eth_units.GweiPerEth
	return artemis_eth_units.GweiMultiple(gweiInEther / 4)
}

func (a *AuxiliaryTradingUtils) isProfitHigherThanGasFee(tf *web3_client.TradeExecutionFlow) (bool, error) {
	if tf.FrontRunTrade.TotalGasCost == 0 {
		return false, errors.New("front run gas cost is 0")
	}
	if tf.SandwichTrade.TotalGasCost == 0 {
		return false, errors.New("sandwich gas cost is 0")
	}
	totalGasCost := tf.FrontRunTrade.TotalGasCost + tf.SandwichTrade.TotalGasCost
	if !artemis_eth_units.IsXGreaterThanY(tf.SandwichTrade.AmountOut, artemis_eth_units.NewBigIntFromUint(totalGasCost)) {
		return false, errors.New("profit is not higher than gas fee")
	}
	return true, nil
}

func (a *AuxiliaryTradingUtils) isTradingEnabledOnToken(tk string) (bool, error) {
	tan := artemis_trading_cache.TokenMap[tk].TradingEnabled
	if tan == nil {
		return false, errors.New("token not found in cache")
	}
	if *tan {
		return *tan, nil
	} else {
		return false, errors.New("trading is disabled")
	}
}

// in sandwich trade the tokenIn on the first trade is the profit currency
func (a *AuxiliaryTradingUtils) isProfitTokenAcceptable(ctx context.Context, tf *web3_client.TradeExecutionFlow) (bool, error) {
	wethAddr := a.getChainSpecificWETH()
	if tf.FrontRunTrade.AmountInAddr.String() != wethAddr.String() {
		return false, errors.New("profit token is not WETH")
	}
	ok, err := a.isTradingEnabledOnToken(tf.FrontRunTrade.AmountOutAddr.String())
	if err != nil {
		log.Info().Interface("tf.FrontRunTrade.AmountOutAddr", tf.FrontRunTrade.AmountOutAddr.String()).Msg("trading is disabled for token")
		return false, err
	}
	err = tf.FrontRunTrade.GetGasUsageForAllTxs(ctx, a.U.Web3Client.Web3Actions)
	if err != nil {
		return false, errors.New("could not get gas usage for front run txs")
	}
	err = tf.SandwichTrade.GetGasUsageForAllTxs(ctx, a.U.Web3Client.Web3Actions)
	if err != nil {
		return false, errors.New("could not get gas usage for sandwich txs")
	}
	ok, err = a.isProfitHigherThanGasFee(tf)
	if err != nil {
		return false, err
	}
	if !ok {
		log.Info().Interface("tf.SandwichTrade.AmountOut", tf.SandwichTrade.AmountOut).Msg("profit is not higher than gas fee")
		return false, errors.New("profit is not higher than gas fee")
	}
	if artemis_eth_units.IsXGreaterThanY(tf.FrontRunTrade.AmountIn, a.maxTradeSize()) {
		log.Info().Interface("tf.FrontRunTrade.AmountIn", tf.FrontRunTrade.AmountIn).Interface("maxTradeSize", a.maxTradeSize()).Msg("trade size is higher than max trade size")
		return false, errors.New("trade size is higher than max trade size")
	}
	return true, nil
}

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

// 0.333 WETH at the moment
// minWethAmountGwei := 330000000
func maxTradeSize() *big.Int {
	gweiInEther := artemis_eth_units.GweiPerEth
	return artemis_eth_units.GweiMultiple(gweiInEther / 3)
}

func isProfitHigherThanGasFee(tf *web3_client.TradeExecutionFlow) (bool, error) {
	log.Info().Msgf("isProfitHigherThanGasFee: front run gas cost: %d", tf.FrontRunTrade.TotalGasCost)
	if tf.FrontRunTrade.TotalGasCost == 0 {
		return false, errors.New("front run gas cost is 0")
	}
	log.Info().Msgf("isProfitHigherThanGasFee: sandwich run gas cost: %d", tf.SandwichTrade.TotalGasCost)
	if tf.SandwichTrade.TotalGasCost == 0 {
		return false, errors.New("sandwich gas cost is 0")
	}
	totalGasCost := tf.FrontRunTrade.TotalGasCost + tf.SandwichTrade.TotalGasCost
	log.Info().Msgf("isProfitHigherThanGasFee: totalGasCost: %d", totalGasCost)
	log.Info().Msgf("tf.SandwichTrade.AmountOut: %s", tf.SandwichTrade.AmountOut.String())
	if !artemis_eth_units.IsXGreaterThanY(tf.SandwichTrade.AmountOut, artemis_eth_units.NewBigIntFromUint(totalGasCost)) {
		return false, errors.New("profit is not higher than gas fee")
	}
	return true, nil
}

func isBundleProfitHigherThanGasFee(bundle *MevTxGroup, tf *web3_client.TradeExecutionFlow) (bool, error) {
	totalGasCost := bundle.TotalGasCost
	log.Info().Msgf("isProfitHigherThanGasFee: expectedProfit: %d", tf.SandwichPrediction.ExpectedProfit)
	log.Info().Msgf("isProfitHigherThanGasFee: totalGasCost: %d", totalGasCost)

	if artemis_eth_units.IsXGreaterThanY(tf.SandwichPrediction.ExpectedProfit, totalGasCost) {
		return false, errors.New("profit is not higher than gas fee")
	}
	return true, nil
}

func IsTradingEnabledOnToken(tk string) (bool, error) {
	tan := artemis_trading_cache.TokenMap[tk].TradingEnabled
	if tan == nil {
		log.Warn().Str("token", tk).Msg("IsTradingEnabledOnToken: token not found in cache")
		return false, errors.New("IsTradingEnabledOnToken: token not found in cache")
	}
	if *tan {
		return *tan, nil
	} else {
		log.Warn().Str("token", tk).Msg("IsTradingEnabledOnToken: token trading is disabled")
		return false, errors.New("IsTradingEnabledOnToken: trading is disabled")
	}
}

// IsProfitTokenAcceptable in sandwich trade the tokenIn on the first trade is the profit currency
func IsProfitTokenAcceptable(ctx context.Context, w3c web3_client.Web3Client, tf *web3_client.TradeExecutionFlow) (bool, error) {
	log.Info().Str("txHash", tf.Tx.Hash().String()).Interface("tf.FrontRunTrade.AmountInAddr.String() ", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr.String()", tf.FrontRunTrade.AmountOutAddr.String()).Msg("is profit token acceptable")
	// just assumes mainnet for now
	if tf.FrontRunTrade.AmountInAddr.String() == tf.FrontRunTrade.AmountOutAddr.String() {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Interface("tf.FrontRunTrade.AmountInAddr.String() ", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr.String()", tf.FrontRunTrade.AmountOutAddr.String()).Msg("profit token is not WETH or ETH")
		return false, errors.New("tokenIn and tokenOut are the same")
	}
	//wethAddr := artemis_trading_constants.WETH9ContractAddress
	//if tf.FrontRunTrade.AmountInAddr.String() != wethAddr && tf.FrontRunTrade.AmountInAddr.String() != artemis_trading_constants.ZeroAddress {
	//	log.Warn().Str("txHash", tf.Tx.Hash().String()).Interface("tf.FrontRunTrade.AmountInAddr.String() ", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr.String()", tf.FrontRunTrade.AmountOutAddr.String()).Msg("profit token is not WETH or ETH")
	//	return false, errors.New("profit token is not WETH or ETH")
	//}

	ok, err := IsTradingEnabledOnToken(tf.FrontRunTrade.AmountOutAddr.String())
	if err != nil {
		log.Info().Interface("tf.FrontRunTrade.AmountInAddr.String()", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr", tf.FrontRunTrade.AmountOutAddr.String()).Msg("trading is disabled for token")
		return false, err
	}
	if !ok {
		log.Info().Msg("trading not enabled on token")
		return false, errors.New("trading not enabled on token")
	}
	if artemis_eth_units.IsXGreaterThanY(tf.FrontRunTrade.AmountIn, maxTradeSize()) {
		log.Info().Interface("tf.FrontRunTrade.AmountIn", tf.FrontRunTrade.AmountIn).Interface("maxTradeSize", maxTradeSize()).Msg("trade size is higher than max trade size")
		return false, errors.New("trade size is higher than max trade size")
	}
	// 0.05 ETH at the moment, ~$100
	minEthAmountGwei := 100000000 / 2
	ok, err = CheckEthBalanceGreaterThan(ctx, w3c, artemis_eth_units.GweiMultiple(minEthAmountGwei))
	if err != nil {
		log.Err(err).Msg("could not check eth balance")
		return false, err
	}
	if !ok {
		return false, errors.New("ETH balance is not enough")
	}
	ok, err = CheckMainnetAuxWETHBalanceGreaterThan(ctx, w3c, maxTradeSize())
	if err != nil {
		log.Err(err).Msg("could not check aux weth balance")
		return false, err
	}
	if !ok {
		return false, errors.New("ETH balance is not enough")
	}
	return true, nil
}

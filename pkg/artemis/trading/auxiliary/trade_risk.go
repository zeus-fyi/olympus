package artemis_trading_auxiliary

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
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
	if tf.Tx == nil {
		log.Warn().Msg("IsProfitTokenAcceptable: tx is nil")
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Interface("tf.FrontRunTrade.AmountInAddr.String() ", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr.String()", tf.FrontRunTrade.AmountOutAddr.String()).Msg("IsProfitTokenAcceptable: is profit token acceptable")
	// just assumes mainnet for now
	if tf.FrontRunTrade.AmountInAddr.String() == tf.FrontRunTrade.AmountOutAddr.String() {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Interface("tf.FrontRunTrade.AmountInAddr.String() ", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr.String()", tf.FrontRunTrade.AmountOutAddr.String()).Msg("IsProfitTokenAcceptable: profit token is not WETH or ETH")
		return false, errors.New("IsProfitTokenAcceptable: tokenIn and tokenOut are the same")
	}
	//wethAddr := artemis_trading_constants.WETH9ContractAddress
	if tf.FrontRunTrade.AmountInAddr.String() == artemis_trading_constants.ZeroAddress && tf.FrontRunTrade.AmountOutAddr.String() == artemis_trading_constants.ZeroAddress {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Interface("tf.FrontRunTrade.AmountInAddr.String() ", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr.String()", tf.FrontRunTrade.AmountOutAddr.String()).Msg("IsProfitTokenAcceptable: empty token addresses")
		return false, errors.New("IsProfitTokenAcceptable: profit token addresses are empty")
	}
	if tf.SandwichTrade.AmountInAddr.String() == artemis_trading_constants.ZeroAddress && tf.SandwichTrade.AmountOutAddr.String() == artemis_trading_constants.ZeroAddress {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Interface("tf.SandwichTrade.AmountInAddr.String() ", tf.SandwichTrade.AmountInAddr.String()).Interface("tf.SandwichTrade.AmountOutAddr.String()", tf.SandwichTrade.AmountOutAddr.String()).Msg("IsProfitTokenAcceptable: empty token addresses")
		return false, errors.New("IsProfitTokenAcceptable: profit token addresses are empty")
	}
	if tf.SandwichTrade.AmountOutAddr.String() != artemis_trading_constants.WETH9ContractAddress {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Interface("tf.SandwichTrade.AmountOutAddr.String()", tf.SandwichTrade.AmountOutAddr.String()).Interface("tf.FrontRunTrade.AmountInAddr.String() ", tf.FrontRunTrade.AmountInAddr.String()).Msg("IsProfitTokenAcceptable: profit token is not the same")
		return false, fmt.Errorf("IsProfitTokenAcceptable: profit token is not WETH %s", tf.SandwichTrade.AmountOutAddr.String())
	}

	log.Info().Str("txHash", tf.Tx.Hash().String()).Interface("tf.FrontRunTrade.AmountInAddr.String() ", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr.String()", tf.FrontRunTrade.AmountOutAddr.String()).Msg("IsProfitTokenAcceptable: profit token is acceptable")
	ok1 := artemis_eth_units.IsStrXLessThanEqZeroOrOne(tf.FrontRunTrade.AmountIn.String())
	ok2 := artemis_eth_units.IsStrXLessThanEqZeroOrOne(tf.FrontRunTrade.AmountOut.String())
	ok3 := artemis_eth_units.IsStrXLessThanEqZeroOrOne(tf.UserTrade.AmountIn.String())
	ok4 := artemis_eth_units.IsStrXLessThanEqZeroOrOne(tf.UserTrade.AmountOut.String())
	ok5 := artemis_eth_units.IsStrXLessThanEqZeroOrOne(tf.SandwichTrade.AmountIn.String())
	ok6 := artemis_eth_units.IsStrXLessThanEqZeroOrOne(tf.SandwichTrade.AmountOut.String())
	ok7 := artemis_eth_units.IsStrXLessThanEqZeroOrOne(tf.SandwichPrediction.ExpectedProfit.String())
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		log.Warn().Msg("IsProfitTokenAcceptable: one of the trade amountsIn or amountsOut is zero")
		return false, errors.New("one of the trade amountsIn or amountsOut is zero")
	}

	log.Info().Str("txHash", tf.Tx.Hash().String()).Interface("tf.FrontRunTrade.AmountInAddr.String() ", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr.String()", tf.FrontRunTrade.AmountOutAddr.String()).Msg("IsProfitTokenAcceptable: profit amount is not zero")
	ok, err := IsTradingEnabledOnToken(tf.FrontRunTrade.AmountOutAddr.String())
	if err != nil {
		log.Info().Interface("tf.FrontRunTrade.AmountInAddr.String()", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr", tf.FrontRunTrade.AmountOutAddr.String()).Msg("IsProfitTokenAcceptable: trading is disabled for token")
		return false, err
	}
	if !ok {
		log.Info().Msg("IsProfitTokenAcceptable: trading not enabled on token")
		return false, errors.New("IsProfitTokenAcceptable: trading not enabled on token")
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Interface("tf.FrontRunTrade.AmountInAddr.String() ", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountOutAddr.String()", tf.FrontRunTrade.AmountOutAddr.String()).Msg("IsProfitTokenAcceptable: trading token is enabled")

	if artemis_eth_units.IsXGreaterThanY(tf.FrontRunTrade.AmountIn, maxTradeSize()) {
		log.Info().Str("tf.FrontRunTrade.AmountInAddr", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountIn", tf.FrontRunTrade.AmountIn).Interface("maxTradeSize", maxTradeSize()).Msg("IsProfitTokenAcceptable: trade size is higher than max trade size")
		return false, errors.New("IsProfitTokenAcceptable: trade size is higher than max trade size")
	}
	// 0.05 ETH at the moment, ~$100
	minEthAmountGwei := 100000000 / 2
	ok, err = CheckEthBalanceGreaterThan(ctx, w3c, artemis_eth_units.GweiMultiple(minEthAmountGwei))
	if err != nil {
		log.Warn().Err(err).Msg("IsProfitTokenAcceptable: could not check eth balance")
		log.Err(err).Msg("IsProfitTokenAcceptable: could not check eth balance")
		return false, err
	}
	if !ok {
		log.Warn().Msg("IsProfitTokenAcceptable: ETH balance is not enough")
		return false, errors.New("IsProfitTokenAcceptable: ETH balance is not enough")
	}
	ok, err = CheckMainnetAuxWETHBalanceGreaterThan(ctx, w3c, maxTradeSize())
	if err != nil {
		log.Warn().Err(err).Msg("IsProfitTokenAcceptable: could not check aux weth balance")
		log.Err(err).Msg("IsProfitTokenAcceptable: could not check aux weth balance")
		return false, err
	}
	if !ok {
		return false, errors.New("ETH balance is not enough")
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Str("tradeMethod", tf.Trade.TradeMethod)
	log.Info().Interface("tf.FrontRunTrade.AmountInAddr.String()", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountIn", tf.FrontRunTrade.AmountIn.String()).Interface("tf.FrontRunTrade.AmountOutAddr", tf.FrontRunTrade.AmountOutAddr.String()).Msg("IsProfitTokenAcceptable: profit token is acceptable")
	log.Info().Interface("tf.FrontRunTrade", tf.FrontRunTrade).Msg("IsProfitTokenAcceptable: profit token is acceptable")
	log.Info().Interface("tf.UserTrade", tf.UserTrade).Msg("IsProfitTokenAcceptable: profit token is acceptable")
	log.Info().Interface("tf.SandwichTrade", tf.SandwichTrade).Msg("IsProfitTokenAcceptable: profit token is acceptable")
	return true, nil
}

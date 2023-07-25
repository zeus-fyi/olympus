package artemis_realtime_trading

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

/*
  adding in other filters here
	  - filter by token
	  - filter by profit
	  - filter by risk score
	  - adds sourcing of new blocks
*/

func EntryTxFilter(ctx context.Context, tx *types.Transaction) error {
	if tx.To() == nil {
		return errors.New("dat: EntryTxFilter, tx.To() is nil")
	}
	if len(tx.Hash().String()) <= 0 {
		return errors.New("dat: EntryTxFilter, tx.Hash().String() is nil")
	}
	return nil
}

func (a *ActiveTrading) SimTxFilter(ctx context.Context, tfSlice []web3_client.TradeExecutionFlowJSON) error {
	var addresses []accounts.Address
	for _, tf := range tfSlice {
		addresses = append(addresses, tf.UserTrade.AmountInAddr)
		addresses = append(addresses, tf.UserTrade.AmountOutAddr)
	}
	for _, addr := range addresses {
		if artemis_trading_cache.TokenMap[addr.String()].BalanceOfSlotNum < 0 {
			return errors.New("SimTxFilter: EntryTxFilter, balanceOf not cracked yet")
		}
		num := artemis_trading_cache.TokenMap[addr.String()].TransferTaxNumerator
		den := artemis_trading_cache.TokenMap[addr.String()].TransferTaxDenominator
		if num == nil || den == nil {
			return errors.New("SimTxFilter: EntryTxFilter, transfer tax not set")
		}
		if *num == 0 || *den == 0 {
			return errors.New("SimTxFilter: EntryTxFilter, transfer tax not set")
		}
	}
	return nil
}

func ActiveTradingFilterSlice(ctx context.Context, w3c web3_client.Web3Client, tf []web3_client.TradeExecutionFlowJSON) error {
	for _, tradeFlow := range tf {
		tfInt, err := tradeFlow.ConvertToBigIntType()
		if err != nil {
			return err
		}
		err = ActiveTradingFilter(ctx, w3c, tfInt, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func ActiveTradeMethodFilter(ctx context.Context, tm string, m *metrics_trading.TradingMetrics) error {
	switch tm {
	case artemis_trading_constants.SwapExactETHForTokens:
	case artemis_trading_constants.SwapTokensForExactETH:
	case artemis_trading_constants.SwapTokensForExactTokens:
	case artemis_trading_constants.SwapExactTokensForETH:
	case artemis_trading_constants.SwapExactTokensForTokens:
	case artemis_trading_constants.SwapETHForExactTokens:
	case swapExactInputSingle, swapExactOutputSingle:
	case exactInput, exactOutput: //"exactInputSingle", "exactOutputSingle":
	case artemis_trading_constants.SwapExactTokensForETHSupportingFeeOnTransferTokens:
	case artemis_trading_constants.SwapExactETHForTokensSupportingFeeOnTransferTokens:
	case artemis_trading_constants.SwapExactTokensForTokensSupportingFeeOnTransferTokens:
	case artemis_trading_constants.Multicall, artemis_trading_constants.Execute0, artemis_trading_constants.Execute:
	case artemis_trading_constants.V2SwapExactIn, artemis_trading_constants.V2SwapExactOut:
	case artemis_trading_constants.V3SwapExactIn, artemis_trading_constants.V3SwapExactOut:
	default:
		log.Warn().Str("tf.Trade.TradeMethod", tm).Msg("dat: ActiveTradingFilter: method not supported for now")
		return fmt.Errorf("dat: ActiveTradingFilter: %s method not supported for now", tm)
	}
	return nil
}
func ActiveTradingFilter(ctx context.Context, w3c web3_client.Web3Client, tf web3_client.TradeExecutionFlow, m *metrics_trading.TradingMetrics) error {
	err := ActiveTradeMethodFilter(ctx, tf.Trade.TradeMethod, m)
	if err != nil {
		return err
	}

	acceptable, err := artemis_trading_auxiliary.IsProfitTokenAcceptable(&tf, m)
	if err != nil {
		log.Err(err).Msg("ActiveTradingFilter: profit token not acceptable")
		return err
	}
	if !acceptable {
		return errors.New("ActiveTradingFilter: profit token not acceptable")
	}

	// 0.012 eth ~$22
	okProfitAmount := artemis_eth_units.IsXLessThanY(tf.SandwichPrediction.ExpectedProfit, artemis_eth_units.GweiMultiple(12000000))
	if okProfitAmount {
		log.Warn().Str("tf.Tx.Hash", tf.Tx.Hash().String()).Interface("tf.SandwichPrediction", tf.SandwichPrediction).Msg("ActiveTradingFilter: SandwichPrediction profit amount not sufficient")
		return errors.New("tf.SandwichPrediction: one of the trade amountsIn or amountsOut is zero")
	}

	if artemis_eth_units.IsXGreaterThanY(tf.FrontRunTrade.AmountIn, artemis_trading_auxiliary.MaxTradeSize()) {
		if m != nil {
			m.ErrTrackingMetrics.CountTradeSizeErr()
		}
		log.Info().Str("tf.FrontRunTrade.AmountInAddr", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountIn", tf.FrontRunTrade.AmountIn).Interface("maxTradeSize", artemis_trading_auxiliary.MaxTradeSize()).Msg("IsProfitTokenAcceptable: trade size is higher than max trade size")
		return errors.New("IsProfitTokenAcceptable: trade size is higher than max trade size")
	}

	if m != nil {
		m.StageProgressionMetrics.CountCheckpointOneMarker()
	}
	// 0.025 ETH at the moment, ~$50
	minEthAmountGwei := 50000000 / 2
	ok, err := artemis_trading_auxiliary.CheckEthBalanceGreaterThan(context.Background(), w3c, artemis_eth_units.GweiMultiple(minEthAmountGwei))
	if err != nil {
		log.Warn().Err(err).Msg("IsProfitTokenAcceptable: could not check eth balance")
		log.Err(err).Msg("IsProfitTokenAcceptable: could not check eth balance")
		return err
	}
	if !ok {
		log.Warn().Msg("IsProfitTokenAcceptable: ETH balance is not enough")
		return errors.New("IsProfitTokenAcceptable: ETH balance is not enough")
	}

	ok, err = artemis_trading_auxiliary.CheckMainnetAuxWETHBalanceGreaterThan(context.Background(), w3c, artemis_trading_auxiliary.MaxTradeSize())
	if err != nil {
		log.Warn().Err(err).Msg("IsProfitTokenAcceptable: could not check aux weth balance")
		log.Err(err).Msg("IsProfitTokenAcceptable: could not check aux weth balance")
		return err
	}
	if !ok {
		log.Warn().Msg("IsProfitTokenAcceptable: WETH balance is not enough")
		return errors.New("IsProfitTokenAcceptable: WETH balance is not enough")
	}
	if m != nil {
		m.StageProgressionMetrics.CountCheckpointTwoMarker()
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Str("tradeMethod", tf.Trade.TradeMethod).Interface("tf.FrontRunTrade.AmountInAddr.String()", tf.FrontRunTrade.AmountInAddr.String()).Interface("tf.FrontRunTrade.AmountIn", tf.FrontRunTrade.AmountIn.String()).Interface("tf.FrontRunTrade.AmountOutAddr", tf.FrontRunTrade.AmountOutAddr.String()).Msg("IsProfitTokenAcceptable: profit token is acceptable")
	log.Info().Str("txHash", tf.Tx.Hash().String()).Str("tradeMethod", tf.Trade.TradeMethod).Interface("tf.FrontRunTrade", tf.FrontRunTrade).Msg("IsProfitTokenAcceptable: profit token is acceptable")
	log.Info().Str("txHash", tf.Tx.Hash().String()).Str("tradeMethod", tf.Trade.TradeMethod).Interface("tf.UserTrade", tf.UserTrade).Msg("IsProfitTokenAcceptable: profit token is acceptable")
	log.Info().Str("txHash", tf.Tx.Hash().String()).Str("tradeMethod", tf.Trade.TradeMethod).Interface("tf.SandwichTrade", tf.SandwichTrade).Msg("IsProfitTokenAcceptable: profit token is acceptable")

	if artemis_eth_units.IsXLessThanY(tf.SandwichPrediction.ExpectedProfit, tf.FrontRunTrade.AmountIn) {
		log.Warn().Interface("tf.SandwichPrediction.ExpectedProfit", tf.SandwichPrediction.ExpectedProfit).Interface(" tf.FrontRunTrade.AmountIn", tf.FrontRunTrade.AmountIn).Msg("ActiveTradingFilter: profit less than trade amount in")
		return fmt.Errorf("dat: ActiveTradingFilter: profit margin min")
	}

	//ok, err := a.GetAuxClient().IsTradingEnabledOnToken(tf.UserTrade.AmountOutAddr.String())
	//if err != nil {
	//	log.Err(err).Msg("dat: ActiveTradingFilter: trading not enabled for token")
	//	return err
	//}
	//if !ok {
	//	return fmt.Errorf("dat: ActiveTradingFilter: trading not enabled for token")
	//}
	log.Info().Interface("tf.FrontRunTrade", tf.FrontRunTrade).Msg("ActiveTradingFilter: passed")
	log.Info().Interface("tf.UserTrade", tf.UserTrade).Msg("ActiveTradingFilter: passed")
	log.Info().Interface("tf.SandwichTrade", tf.SandwichTrade).Msg("ActiveTradingFilter: passed")
	log.Info().Interface("tf.Tx.Hash", tf.Tx.Hash()).Msg("ActiveTradingFilter: passed")

	return nil
}

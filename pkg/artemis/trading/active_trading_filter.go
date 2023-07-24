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

	if m != nil {
		m.StageProgressionMetrics.CountCheckpointOneMarker()
	}

	_, err = artemis_trading_auxiliary.IsProfitTokenAcceptable(ctx, w3c, &tf, m)
	if err != nil {
		log.Err(err).Msg("ActiveTradingFilter: profit token not acceptable")
		return err
	}

	if m != nil {
		m.StageProgressionMetrics.CountCheckpointTwoMarker()
	}

	// ~$15
	profitMin := artemis_eth_units.GweiMultiple(10000000)
	if artemis_eth_units.IsXLessThanY(tf.SandwichPrediction.ExpectedProfit, profitMin) {
		log.Warn().Interface("tf.SandwichPrediction.ExpectedProfit", tf.SandwichPrediction.ExpectedProfit).Interface("profitMarginMin", profitMin).Msg("ActiveTradingFilter: profit margin min")
		return fmt.Errorf("dat: ActiveTradingFilter: profit margin min")
	}

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

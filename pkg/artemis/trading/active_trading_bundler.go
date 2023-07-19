package artemis_realtime_trading

import (
	"context"

	"github.com/rs/zerolog/log"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) ProcessBundleStage(ctx context.Context, w3c web3_client.Web3Client, tfSlice []web3_client.TradeExecutionFlowJSON, m *metrics_trading.TradingMetrics) error {
	for _, tradeFlow := range tfSlice {
		tf := tradeFlow.ConvertToBigIntType()
		err := ActiveTradingFilter(ctx, w3c, tf)
		if err != nil {
			log.Err(err).Msg("ProcessBundleStage: failed to pass active filter trade")
			err = nil
			continue
		}
		log.Info().Msgf("ProcessBundleStage: passed active filter trade: %s", tf.Tx.Hash().String())
		m.StageProgressionMetrics.CountPostActiveTradingFilter(1)
		resp, _, err := a.GetAuxClient().StagingPackageSandwichAndCall(ctx, &tf)
		if err != nil {
			log.Err(err).Msg("ProcessBundleStage: failed to package sandwich")
			err = nil
			continue
		}
		log.Info().Msgf("ProcessBundleStage: StagingPackageSandwichAndCall passed bundle hash: %s", resp.BundleHash)
		if resp != nil {
			log.Info().Interface("fbCallResp", resp).Msg("sent sandwich")
			m.StageProgressionMetrics.CountSentFlashbotsBundleSubmission(1)
		}
	}
	return nil
}

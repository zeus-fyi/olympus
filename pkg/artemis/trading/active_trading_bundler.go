package artemis_realtime_trading

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) ProcessBundleStage(ctx context.Context, bn uint64, tfSlice []web3_client.TradeExecutionFlowJSON) error {
	for _, tradeFlow := range tfSlice {
		tf := tradeFlow.ConvertToBigIntType()
		_, err := a.a.PackageSandwich(ctx, &tf)
		if err != nil {
			log.Err(err).Msg("failed to package sandwich")
			return err
		}
	}
	return nil
}

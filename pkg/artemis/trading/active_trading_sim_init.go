package artemis_realtime_trading

import (
	"context"

	"github.com/rs/zerolog/log"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) InitActiveTradingSimEnv(ctx context.Context, tf *web3_client.TradeExecutionFlow) error {
	approveTx, err := a.a.U.ApproveSpender(ctx, artemis_trading_constants.WETH9ContractAddress, artemis_trading_constants.Permit2SmartContractAddress, artemis_eth_units.MaxUINT)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	secondToken := tf.FrontRunTrade.AmountInAddr.String()
	if tf.FrontRunTrade.AmountInAddr.String() == artemis_trading_constants.WETH9ContractAddress {
		secondToken = tf.FrontRunTrade.AmountOutAddr.String()
	}
	approveTx, err = a.a.U.ApproveSpender(ctx, secondToken, artemis_trading_constants.Permit2SmartContractAddress, artemis_eth_units.MaxUINT)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	return nil
}

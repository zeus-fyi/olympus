package artemis_realtime_trading

import (
	"context"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/metachris/flashbotsrpc"
	"github.com/rs/zerolog/log"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) SendToBundleStack(ctx context.Context, tx *types.Transaction) {
}

func (a *ActiveTrading) BundleTxs(ctx context.Context, tfSlice []web3_client.TradeExecutionFlowJSON) ([]artemis_flashbots.MevTxBundle, error) {
	var bundles []artemis_flashbots.MevTxBundle
	for _, tradeFlow := range tfSlice {
		tf := tradeFlow.ConvertToBigIntType()
		// todo, shouldn't necessarily bypass sim stage
		err := a.SimToPackageTxBundle(ctx, &tf, true)
		if err != nil {
			return nil, err
		}
		if tf.Bundle != nil {
			bundles = append(bundles, *tf.Bundle)
			// todo update metric here
		}
	}
	return bundles, nil
}

func (a *ActiveTrading) SubmitCallBundle(ctx context.Context, bn uint64, bundles []artemis_flashbots.MevTxBundle) error {
	for _, bundle := range bundles {
		param := flashbotsrpc.FlashbotsCallBundleParam{
			BlockNumber: "0x" + strconv.FormatUint(bn+2, 10),
			Txs:         bundle.Txs,
		}
		resp, ferr := a.a.CallBundle(ctx, param)
		if ferr != nil {
			log.Err(ferr).Msg("failed to send flashbots bundle")
			return ferr
		}
		log.Info().Msgf("Flashbots bundle sent, resp: %v", resp)
	}
	return nil
}

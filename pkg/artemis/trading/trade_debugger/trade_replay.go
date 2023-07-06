package artemis_trade_debugger

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (t *TradeDebugger) Replayer(ctx context.Context, txHash string) error {
	hist, err := t.lookupMevTx(ctx, txHash)
	if err != nil {
		return err
	}
	tfSteps := make([]web3_client.TradeExecutionFlow, 0)
	for _, h := range hist {
		tfSteps = append(tfSteps, h.TradePrediction)
	}
	for _, tf := range tfSteps {
		fmt.Println(tf.CurrentBlockNumber.String())
		err = t.SetupCleanEnvironment(ctx, &tf)
		if err != nil {
			return err
		}

	}
	return nil
}

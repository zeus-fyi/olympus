package artemis_trade_debugger

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (t *TradeDebugger) Replay(ctx context.Context, tf *web3_client.TradeExecutionFlow) error {
	err := t.ResetAndSetupPreconditions(ctx, tf)
	if err != nil {
		return err
	}
	return nil
}

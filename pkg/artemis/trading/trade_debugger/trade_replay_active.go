package artemis_trade_debugger

import (
	"context"

	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
)

func (t *TradeDebugger) ReplayUsingActiveTradeFlow(ctx context.Context, txHash string, fromMempoolTx bool) error {
	artemis_realtime_trading.IngestTx(ctx, *t.dat.SimW3c(), nil, nil)
	return nil
}

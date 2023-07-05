package artemis_trade_debugger

import artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"

type TradeDebugger struct {
	ActiveTrading artemis_realtime_trading.ActiveTrading
}

func NewTradeDebugger(a artemis_realtime_trading.ActiveTrading) TradeDebugger {
	return TradeDebugger{
		ActiveTrading: a,
	}
}

func (t *TradeDebugger) GetTxFromHash() {
}

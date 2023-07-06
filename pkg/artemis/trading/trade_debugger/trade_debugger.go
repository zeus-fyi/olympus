package artemis_trade_debugger

import (
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type TradeDebugger struct {
	UniswapClient *web3_client.UniswapClient
	ActiveTrading artemis_realtime_trading.ActiveTrading
}

func NewTradeDebugger(a artemis_realtime_trading.ActiveTrading, u *web3_client.UniswapClient) TradeDebugger {
	return TradeDebugger{
		ActiveTrading: a,
		UniswapClient: u,
	}
}

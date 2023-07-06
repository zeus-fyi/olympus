package artemis_trade_debugger

import (
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
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

type HistoricalAnalysisDebug struct {
	HistoricalAnalysis artemis_mev_models.HistoricalAnalysis
	TradePrediction    web3_client.TradeExecutionFlow
}

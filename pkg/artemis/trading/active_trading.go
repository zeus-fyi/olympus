package artemis_realtime_trading

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type ActiveTrading struct {
	u *web3_client.UniswapClient
	m metrics_trading.TradingMetrics
}

func NewActiveTradingModule(u *web3_client.UniswapClient) ActiveTrading {
	return ActiveTrading{u, metrics_trading.NewTradingMetrics()}
}

func (a *ActiveTrading) IngestTx(ctx context.Context, tx *types.Transaction) {
	a.FilterTx(ctx, tx)
	a.DecodeTx(ctx, tx)
	a.ProcessTx(ctx, tx)
	a.SimulateTx(ctx, tx)
	a.SendToBundleStack(ctx, tx)
}

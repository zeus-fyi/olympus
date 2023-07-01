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

func NewActiveTradingModule(u *web3_client.UniswapClient, tm metrics_trading.TradingMetrics) ActiveTrading {
	ctx := context.Background()
	InitTokenFilter(ctx)
	return ActiveTrading{u, tm}
}

func (a *ActiveTrading) IngestTx(ctx context.Context, tx *types.Transaction) {
	// TODO add metrics pass rate & timing for each stage
	err := a.EntryTxFilter(ctx, tx)
	if err != nil {
		return
	}
	err = a.DecodeTx(ctx, tx)
	if err != nil {
		return
	}
	err = a.ProcessTxs(ctx)
	if err != nil {
		return
	}
	err = a.SimTxFilter(ctx, tx)
	if err != nil {
		return
	}
	err = a.SimulateTx(ctx, tx)
	if err != nil {
		return
	}
	//a.SendToBundleStack(ctx, tx)
}

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
	return ActiveTrading{u, tm}
}

func (a *ActiveTrading) IngestTx(ctx context.Context, tx *types.Transaction) error {
	// TODO add metrics pass rate & timing for each stage
	err := a.EntryTxFilter(ctx, tx)
	if err != nil {
		return err
	}
	err = a.DecodeTx(ctx, tx)
	if err != nil {
		return err
	}
	err = a.ProcessTxs(ctx)
	if err != nil {
		return err
	}
	err = a.SimTxFilter(ctx, tx)
	if err != nil {
		return err
	}
	//go func() {
	//	err = a.ProcessTx(ctx, tx)
	//	if err != nil {
	//		return
	//	}
	//}()
	return err
}

func (a *ActiveTrading) ProcessTx(ctx context.Context, tx *types.Transaction) error {
	// TODO, simulate tx needs a clean anvil instance
	err := a.SimulateTx(ctx, tx)
	if err != nil {
		return err
	}
	a.SendToBundleStack(ctx, tx)
	return nil
}

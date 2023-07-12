package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type ActiveTrading struct {
	a *artemis_trading_auxiliary.AuxiliaryTradingUtils
	m metrics_trading.TradingMetrics
}

func (a *ActiveTrading) GetUniswapClient() *web3_client.UniswapClient {
	return a.GetAuxClient().U
}

func (a *ActiveTrading) GetAuxClient() *artemis_trading_auxiliary.AuxiliaryTradingUtils {
	return a.a
}

func NewActiveTradingModuleWithoutMetrics(a *artemis_trading_auxiliary.AuxiliaryTradingUtils) ActiveTrading {
	return ActiveTrading{a: a}
}

func NewActiveTradingModule(a *artemis_trading_auxiliary.AuxiliaryTradingUtils, tm metrics_trading.TradingMetrics) ActiveTrading {
	return ActiveTrading{a, tm}
}

func (a *ActiveTrading) IngestTx(ctx context.Context, tx *types.Transaction) error {
	a.m.StageProgressionMetrics.CountPreEntryFilterTx()
	err := a.EntryTxFilter(ctx, tx)
	if err != nil {
		return err
	}
	a.m.StageProgressionMetrics.CountPostEntryFilterTx()
	mevTxs, err := a.DecodeTx(ctx, tx)
	if err != nil {
		return err
	}
	if len(mevTxs) <= 0 {
		return errors.New("DecodeTx: no txs to process")
	}
	a.m.StageProgressionMetrics.CountPostDecodeTx()
	tfSlice, err := a.ProcessTxs(ctx)
	if err != nil {
		return err
	}
	if len(tfSlice) <= 0 {
		return errors.New("ProcessTxs: no tx flows to simulate")
	}
	a.m.StageProgressionMetrics.CountPostProcessTx(float64(1))
	err = a.SimTxFilter(ctx, tfSlice)
	if err != nil {
		return err
	}
	if len(tfSlice) <= 0 {
		return errors.New("SimTxFilter: no tx flows to simulate")
	}
	a.m.StageProgressionMetrics.CountPostSimFilterTx(float64(1))
	// todo refactor this
	wc := web3_actions.NewWeb3ActionsClient(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL)
	wc.Dial()
	bn, berr := wc.C.BlockNumber(ctx)
	if berr != nil {
		log.Err(berr).Msg("failed to get block number")
		return berr
	}
	wc.Close()
	err = a.SaveMempoolTx(ctx, bn, tfSlice)
	if err != nil {
		return err
	}
	err = a.ProcessBundleStage(ctx, bn, tfSlice)
	if err != nil {
		return err
	}
	return err
}

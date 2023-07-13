package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const (
	//irisBetaSvc = "https://iris.zeus.fyi/v1beta/internal/"
	irisBetaSvc = "http://iris.iris.svc.cluster.local:8080/v1beta/internal/"
)

type ActiveTrading struct {
	a  *artemis_trading_auxiliary.AuxiliaryTradingUtils
	us *ActiveTrading
	m  *metrics_trading.TradingMetrics
}

func (a *ActiveTrading) GetUniswapClient() *web3_client.UniswapClient {
	return a.GetAuxClient().U
}

func (a *ActiveTrading) GetAuxClient() *artemis_trading_auxiliary.AuxiliaryTradingUtils {
	return a.a
}
func (a *ActiveTrading) GetMetricsClient() *metrics_trading.TradingMetrics {
	return a.m
}

func createSimClient() web3_client.Web3Client {
	sw3c := web3_client.NewWeb3ClientFakeSigner(irisBetaSvc)
	sw3c.AddBearerToken(artemis_orchestration_auth.Bearer)
	return sw3c
}
func NewActiveTradingModuleWithoutMetrics(a *artemis_trading_auxiliary.AuxiliaryTradingUtils) ActiveTrading {
	ctx := context.Background()
	us := web3_client.InitUniswapClient(ctx, createSimClient())
	auxSim := artemis_trading_auxiliary.InitAuxiliaryTradingUtils(ctx, us.Web3Client)
	auxSimTrader := ActiveTrading{
		a: &auxSim,
	}
	return ActiveTrading{a: a, us: &auxSimTrader}
}

func NewActiveTradingModule(a *artemis_trading_auxiliary.AuxiliaryTradingUtils, tm *metrics_trading.TradingMetrics) ActiveTrading {
	ctx := context.Background()
	us := web3_client.InitUniswapClient(ctx, createSimClient())
	auxSim := artemis_trading_auxiliary.InitAuxiliaryTradingUtils(ctx, us.Web3Client)
	auxSimTrader := ActiveTrading{
		a: &auxSim,
	}
	at := ActiveTrading{a: a, us: &auxSimTrader, m: tm}
	return at
}

func (a *ActiveTrading) IngestTx(ctx context.Context, tx *types.Transaction) error {
	a.GetMetricsClient().StageProgressionMetrics.CountPreEntryFilterTx()
	err := a.EntryTxFilter(ctx, tx)
	if err != nil {
		return err
	}
	a.GetMetricsClient().StageProgressionMetrics.CountPostEntryFilterTx()
	mevTxs, err := a.DecodeTx(ctx, tx)
	if err != nil {
		return err
	}
	if len(mevTxs) <= 0 {
		return errors.New("DecodeTx: no txs to process")
	}
	a.GetMetricsClient().StageProgressionMetrics.CountPostDecodeTx()
	tfSlice, err := a.ProcessTxs(ctx)
	if err != nil {
		return err
	}
	if len(tfSlice) <= 0 {
		return errors.New("ProcessTxs: no tx flows to simulate")
	}
	a.GetMetricsClient().StageProgressionMetrics.CountPostProcessTx(float64(1))
	err = a.SimTxFilter(ctx, tfSlice)
	if err != nil {
		return err
	}
	if len(tfSlice) <= 0 {
		return errors.New("SimTxFilter: no tx flows to simulate")
	}
	a.GetMetricsClient().StageProgressionMetrics.CountPostSimFilterTx(float64(1))
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
	log.Info().Msg("starting simulation")
	err = a.SimToPackageTxBundles(ctx, tfSlice, false)
	if err != nil {
		return err
	}
	a.GetMetricsClient().StageProgressionMetrics.CountPostSimStage(float64(len(tfSlice)))
	err = a.ActiveTradingFilterSlice(ctx, tfSlice)
	if err != nil {
		return err
	}
	log.Info().Msg("preparing bundles for submission")
	a.GetMetricsClient().StageProgressionMetrics.CountPostActiveTradingFilter(float64(len(tfSlice)))
	err = a.ProcessBundleStage(ctx, tfSlice)
	if err != nil {
		return err
	}
	log.Info().Msg("bundles successfully sent")
	a.GetMetricsClient().StageProgressionMetrics.CountSentFlashbotsBundleSubmission(float64(len(tfSlice)))
	return err
}

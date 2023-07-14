package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const (
	irisBetaSvcExt = "https://iris.zeus.fyi/v1beta/internal/"
	irisBetaSvc    = "http://iris.iris.svc.cluster.local:8080/v1beta/internal/"
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

func createExtSimClient() web3_client.Web3Client {
	sw3c := web3_client.NewWeb3ClientFakeSigner(irisBetaSvcExt)
	sw3c.AddBearerToken(artemis_orchestration_auth.Bearer)
	return sw3c
}

func NewActiveTradingDebugger(usc *web3_client.UniswapClient) ActiveTrading {
	ctx := context.Background()
	auxSim := artemis_trading_auxiliary.InitAuxiliaryTradingUtilsFromUni(ctx, usc)
	auxSimTrader := ActiveTrading{
		a: &auxSim,
	}
	return ActiveTrading{us: &auxSimTrader}
}

func NewActiveTradingModuleWithoutMetrics(a *artemis_trading_auxiliary.AuxiliaryTradingUtils) ActiveTrading {
	ctx := context.Background()
	us := web3_client.InitUniswapClient(ctx, createSimClient())
	us.Web3Client.IsAnvilNode = true
	us.Web3Client.DurableExecution = true
	auxSim := artemis_trading_auxiliary.InitAuxiliaryTradingUtils(ctx, us.Web3Client)
	auxSimTrader := ActiveTrading{
		a: &auxSim,
	}
	return ActiveTrading{a: a, us: &auxSimTrader}
}

func NewActiveTradingModule(a *artemis_trading_auxiliary.AuxiliaryTradingUtils, tm *metrics_trading.TradingMetrics) ActiveTrading {
	ctx := context.Background()
	us := web3_client.InitUniswapClient(ctx, createSimClient())
	us.Web3Client.IsAnvilNode = true
	us.Web3Client.DurableExecution = true
	auxSim := artemis_trading_auxiliary.InitAuxiliaryTradingUtils(ctx, us.Web3Client)
	auxSimTrader := ActiveTrading{
		a: &auxSim,
	}
	at := ActiveTrading{a: a, us: &auxSimTrader, m: tm}

	go artemis_trading_cache.SetActiveTradingBlockCache(ctx)
	return at
}

type ErrWrapper struct {
	Err   error
	Stage string
	Code  int
}

func (a *ActiveTrading) IngestTx(ctx context.Context, tx *types.Transaction) ErrWrapper {
	a.GetMetricsClient().StageProgressionMetrics.CountPreEntryFilterTx()
	err := a.EntryTxFilter(ctx, tx)
	if err != nil {
		return ErrWrapper{Err: err, Stage: "EntryTxFilter"}
	}
	a.GetMetricsClient().StageProgressionMetrics.CountPostEntryFilterTx()
	mevTxs, err := a.DecodeTx(ctx, tx)
	if err != nil {
		return ErrWrapper{Err: err, Stage: "DecodeTx"}
	}
	if len(mevTxs) <= 0 {
		return ErrWrapper{Err: errors.New("DecodeTx: no txs to process"), Stage: "DecodeTx"}
	}
	a.GetMetricsClient().StageProgressionMetrics.CountPostDecodeTx()
	tfSlice, err := a.ProcessTxs(ctx)
	if err != nil {
		log.Err(err).Msg("failed to pass process txs")
		return ErrWrapper{Err: err, Stage: "ProcessTxs"}
	}
	if len(tfSlice) <= 0 {
		return ErrWrapper{Err: errors.New("ProcessTxs: no tx flows to simulate"), Stage: "ProcessTxs"}
	}
	a.GetMetricsClient().StageProgressionMetrics.CountPostProcessTx(float64(1))
	err = a.SimTxFilter(ctx, tfSlice)
	if err != nil {
		log.Err(err).Msg("failed to pass sim tx filter")
		return ErrWrapper{Err: err, Stage: "SimTxFilter"}
	}
	if len(tfSlice) <= 0 {
		return ErrWrapper{Err: errors.New("SimTxFilter: no tx flows to simulate"), Stage: "SimTxFilter"}
	}
	log.Info().Msg("saving mempool tx")
	a.GetMetricsClient().StageProgressionMetrics.CountPostSimFilterTx(float64(1))
	bn, err := artemis_trading_cache.GetLatestBlock(ctx)
	if err != nil {
		log.Err(err).Msg("failed to get latest block")
		return ErrWrapper{Err: err, Stage: "GetLatestBlock"}
	}
	err = a.SaveMempoolTx(ctx, bn, tfSlice)
	if err != nil {
		log.Err(err).Msg("failed to pass sim tx filter")
		return ErrWrapper{Err: err, Stage: "SaveMempoolTx"}
	}
	log.Info().Msg("starting simulation")
	err = a.SimToPackageTxBundles(ctx, tfSlice, false)
	if err != nil {
		log.Err(err).Msg("failed to simulate txs")
		return ErrWrapper{Err: err, Stage: "SimToPackageTxBundles", Code: 200}
	}
	log.Info().Msg("simulation stage complete: starting active trading filter")
	a.GetMetricsClient().StageProgressionMetrics.CountPostSimStage(float64(len(tfSlice)))
	err = a.ActiveTradingFilterSlice(ctx, tfSlice)
	if err != nil {
		log.Err(err).Msg("failed to pass active trading filter")
		return ErrWrapper{Err: err, Stage: "ActiveTradingFilterSlice", Code: 200}
	}
	log.Info().Msg("preparing bundles for submission")
	a.GetMetricsClient().StageProgressionMetrics.CountPostActiveTradingFilter(float64(len(tfSlice)))
	err = a.ProcessBundleStage(ctx, tfSlice)
	if err != nil {
		log.Err(err).Msg("failed to process bundles")
		return ErrWrapper{Err: err, Stage: "ProcessBundleStage", Code: 200}
	}
	log.Info().Msg("bundles successfully sent")
	a.GetMetricsClient().StageProgressionMetrics.CountSentFlashbotsBundleSubmission(float64(len(tfSlice)))
	return ErrWrapper{Err: err, Stage: "Success", Code: 200}
}

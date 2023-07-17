package artemis_realtime_trading

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const (
	irisBetaSvcExt = "https://iris.zeus.fyi/v1beta/internal/"
	irisBetaSvc    = "http://iris.iris.svc.cluster.local/v1beta/internal/"
	irisSvcBeacons = "http://iris.iris.svc.cluster.local/v1beta/internal/router/group?routeGroup=quiknode-mainnet"
)

var (
	CacheBeacon = web3_client.NewWeb3ClientFakeSigner(irisSvcBeacons)
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
	CacheBeacon.AddBearerToken(artemis_orchestration_auth.Bearer)
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

func newActiveTradingModule(a *artemis_trading_auxiliary.AuxiliaryTradingUtils, tm *metrics_trading.TradingMetrics) ActiveTrading {
	ctx := context.Background()
	us := web3_client.InitUniswapClient(ctx, createSimClient())
	us.Web3Client.IsAnvilNode = true
	us.Web3Client.DurableExecution = true
	auxSim := artemis_trading_auxiliary.InitAuxiliaryTradingUtils(ctx, us.Web3Client)
	auxSimTrader := ActiveTrading{
		a: &auxSim,
	}
	at := ActiveTrading{a: a, us: &auxSimTrader, m: tm}

	return at
}
func NewActiveTradingModule(a *artemis_trading_auxiliary.AuxiliaryTradingUtils, tm *metrics_trading.TradingMetrics) ActiveTrading {
	at := newActiveTradingModule(a, tm)
	go artemis_trading_cache.SetActiveTradingBlockCache(context.Background())
	return at
}

type ErrWrapper struct {
	Err   error
	Stage string
	Code  int
}

var txCache = cache.New(time.Hour*24, time.Hour*24)

func (a *ActiveTrading) IngestTx(ctx context.Context, tx *types.Transaction) ErrWrapper {
	at := newActiveTradingModule(a.a, a.m)
	at.GetMetricsClient().StageProgressionMetrics.CountPreEntryFilterTx()
	err := at.EntryTxFilter(ctx, tx)
	if err != nil {
		return ErrWrapper{Err: err, Stage: "EntryTxFilter"}
	}
	at.GetMetricsClient().StageProgressionMetrics.CountPostEntryFilterTx()
	mevTxs, merr := DecodeTx(ctx, tx, a.m)
	if merr != nil {
		log.Err(merr).Msg("decoding txs err")
	}
	if len(mevTxs) <= 0 {
		log.Err(merr).Msg("no mev txs found")
		return ErrWrapper{Err: merr, Stage: "DecodeTx"}
	}
	at.GetMetricsClient().StageProgressionMetrics.CountPostDecodeTx()
	_, err = at.ProcessTxs(ctx, mevTxs)
	if err != nil {
		log.Err(err).Msg("failed to pass process txs")
		return ErrWrapper{Err: err, Stage: "ProcessTxs"}
	}
	return ErrWrapper{Err: err, Stage: "Success", Code: 200}
}

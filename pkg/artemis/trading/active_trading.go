package artemis_realtime_trading

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const (
	irisBetaSvcExt = "https://iris.zeus.fyi/v1beta/internal/"
	irisBetaSvc    = "http://iris.iris.svc.cluster.local/v1beta/internal/"
	irisSvcBeacons = "http://iris.iris.svc.cluster.local/v1beta/internal/router/group?routeGroup=quiknode-mainnet"
)

var (
	TraderClient web3_client.Web3Client
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

var simClient = createSimClient()

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
	createSimClient()
	us := web3_client.InitUniswapClient(ctx, simClient)
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
	createSimClient()
	us := web3_client.InitUniswapClient(ctx, simClient)
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
	traderAcc, err := accounts.CreateAccountFromPkey(a.U.Web3Client.Account.EcdsaPrivateKey())
	if err != nil || traderAcc == nil {
		panic(err)
	}
	TraderClient = web3_client.NewWeb3Client(irisSvcBeacons, traderAcc)
	if len(artemis_orchestration_auth.Bearer) == 0 {
		panic("no bearer token")
	}
	TraderClient.AddBearerToken(artemis_orchestration_auth.Bearer)
	log.Info().Msgf("trader account: %s", traderAcc.Address().String())
	return at
}

type ErrWrapper struct {
	Err   error
	Stage string
	Code  int
}

func IngestTx(ctx context.Context, tx *types.Transaction, m *metrics_trading.TradingMetrics) ErrWrapper {
	m.StageProgressionMetrics.CountPreEntryFilterTx()
	err := EntryTxFilter(ctx, tx)
	if err != nil {
		return ErrWrapper{Err: err, Stage: "EntryTxFilter"}
	}
	m.StageProgressionMetrics.CountPostEntryFilterTx()
	mevTxs, merr := DecodeTx(ctx, tx, m)
	if merr != nil {
		log.Err(merr).Msg("decoding txs err")
		return ErrWrapper{Err: merr, Stage: "DecodeTx"}
	}
	if len(mevTxs) <= 0 {
		log.Err(merr).Msg("no mev txs found")
		return ErrWrapper{Err: merr, Stage: "DecodeTx"}
	}
	m.StageProgressionMetrics.CountPostDecodeTx()

	w3c := web3_client.NewWeb3Client(irisSvcBeacons, TraderClient.Account)
	w3c.AddBearerToken(artemis_orchestration_auth.Bearer)

	log.Info().Msgf("ProcessTxsStage: txs: %d", len(mevTxs))
	tfSlice := ProcessTxs(ctx, &mevTxs, m, w3c.Web3Actions)

	log.Info().Msgf("ProcessBundleStage: txs: %d", len(tfSlice))
	ProcessBundleStage(ctx, w3c, tfSlice, m)
	return ErrWrapper{Err: err, Stage: "Success", Code: 200}
}

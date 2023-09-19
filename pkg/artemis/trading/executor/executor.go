package artemis_trade_executor

import (
	"context"

	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	tyche_metrics "github.com/zeus-fyi/olympus/tyche/metrics"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

var (
	ActiveTrader         artemis_realtime_trading.ActiveTrading
	TradeExecutorMainnet = artemis_trading_auxiliary.AuxiliaryTradingUtils{}
	ActiveGoerliTrader   artemis_realtime_trading.ActiveTrading
	TradeExecutorGoerli  = artemis_trading_auxiliary.AuxiliaryTradingUtils{}
	ActiveTraderW3c      web3_client.Web3Client
)

const irisSvcBeacons = "http://iris.iris.svc.cluster.local/v2/internal/router"

func InitMainnetAuxiliaryTradingUtils(ctx context.Context, age encryption.Age) artemis_trading_auxiliary.AuxiliaryTradingUtils {
	tm := tyche_metrics.TycheMetrics
	acc := InitTradingAccount2(ctx, age)
	wc := web3_client.NewWeb3Client(irisSvcBeacons, &acc)
	wc.AddDefaultEthereumMainnetTableHeader()
	wc.AddBearerToken(artemis_orchestration_auth.Bearer)
	if len(artemis_orchestration_auth.Bearer) == 0 {
		panic("bearer token not set")
	}
	wc.Network = hestia_req_types.Mainnet
	ActiveTraderW3c = wc
	TradeExecutorMainnet = artemis_trading_auxiliary.InitAuxiliaryTradingUtils(ctx, wc)
	if tm == nil {
		ActiveTrader = artemis_realtime_trading.NewActiveTradingModuleWithoutMetrics(&TradeExecutorMainnet)
	} else {
		ActiveTrader = artemis_realtime_trading.NewActiveTradingModule(&TradeExecutorMainnet, &tyche_metrics.TradeMetrics)
	}
	return TradeExecutorMainnet
}

func InitGoerliAuxiliaryTradingUtils(ctx context.Context, age encryption.Age) artemis_trading_auxiliary.AuxiliaryTradingUtils {
	acc := InitTradingAccount(ctx, age)
	cfg := artemis_network_cfgs.ArtemisEthereumGoerli
	wc := web3_client.NewWeb3Client(cfg.NodeURL, &acc)
	wc.Network = cfg.Network
	TradeExecutorGoerli = artemis_trading_auxiliary.InitAuxiliaryTradingUtils(ctx, wc)
	return TradeExecutorGoerli
}

// InitTradingAccount pubkey 0x000025e60C7ff32a3470be7FE3ed1666b0E326e2
func InitTradingAccount(ctx context.Context, age encryption.Age) accounts.Account {
	p := filepaths.Path{
		DirIn:  "keygen",
		DirOut: "keygen",
		FnIn:   "key-4.txt.age",
	}
	r, err := dynamic_secrets.ReadAddress(ctx, p, athena.AthenaS3Manager, age)
	if err != nil {
		panic(err)
	}
	acc, err := dynamic_secrets.GetAccount(r)
	if err != nil {
		panic(err)
	}
	return acc
}

// InitTradingAccount2 pubkey
func InitTradingAccount2(ctx context.Context, age encryption.Age) accounts.Account {
	p := filepaths.Path{
		DirIn:  "keygen",
		DirOut: "keygen",
		FnIn:   "key-6.txt.age",
	}
	r, err := dynamic_secrets.ReadAddress(ctx, p, athena.AthenaS3Manager, age)
	if err != nil {
		panic(err)
	}
	acc, err := dynamic_secrets.GetAccount(r)
	if err != nil {
		panic(err)
	}
	return acc
}

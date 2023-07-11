package artemis_trade_executor

import (
	"context"

	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

var (
	ActiveTrader         artemis_realtime_trading.ActiveTrading
	TradeExecutorMainnet = artemis_trading_auxiliary.AuxiliaryTradingUtils{}
	ActiveGoerliTrader   artemis_realtime_trading.ActiveTrading
	TradeExecutorGoerli  = artemis_trading_auxiliary.AuxiliaryTradingUtils{}
)

func InitMainnetAuxiliaryTradingUtils(ctx context.Context, age encryption.Age) artemis_trading_auxiliary.AuxiliaryTradingUtils {
	acc := InitTradingAccount2(ctx, age)
	cfg := artemis_network_cfgs.ArtemisEthereumMainnet
	TradeExecutorMainnet = artemis_trading_auxiliary.InitAuxiliaryTradingUtils(ctx, cfg.NodeURL, cfg.Network, acc)
	ActiveTrader = artemis_realtime_trading.NewActiveTradingModuleWithoutMetrics(&TradeExecutorMainnet)
	return TradeExecutorMainnet
}

func InitGoerliAuxiliaryTradingUtils(ctx context.Context, age encryption.Age) artemis_trading_auxiliary.AuxiliaryTradingUtils {
	acc := InitTradingAccount(ctx, age)
	cfg := artemis_network_cfgs.ArtemisEthereumGoerli
	TradeExecutorGoerli = artemis_trading_auxiliary.InitAuxiliaryTradingUtils(ctx, cfg.NodeURL, cfg.Network, acc)
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

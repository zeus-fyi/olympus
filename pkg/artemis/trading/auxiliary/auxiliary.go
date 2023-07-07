package artemis_trading_auxiliary

import (
	"context"

	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

var TradingAuxiliary AuxiliaryTradingUtils

func InitAuxiliaryTradingUtils(ctx context.Context, nodeURL string, acc accounts.Account) AuxiliaryTradingUtils {
	TradingAuxiliary = AuxiliaryTradingUtils{
		Web3Client: web3_client.NewWeb3Client(nodeURL, &acc),
	}
	return TradingAuxiliary
}

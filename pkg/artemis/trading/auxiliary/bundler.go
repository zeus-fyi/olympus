package artemis_trading_auxiliary

import (
	"github.com/metachris/flashbotsrpc"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) CreateFlashbotsBundle(ur *web3_client.UniversalRouterExecCmd, bn string) artemis_flashbots.MevTxBundle {
	maxTime := ur.Deadline.Uint64()
	mevBundle := artemis_flashbots.MevTxBundle{
		FlashbotsSendBundleRequest: flashbotsrpc.FlashbotsSendBundleRequest{
			Txs:          []string{},
			BlockNumber:  bn,
			MaxTimestamp: &maxTime,
		},
	}
	mevBundle.AddTxs(a.OrderedTxs...)
	return mevBundle
}

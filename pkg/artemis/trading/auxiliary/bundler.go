package artemis_trading_auxiliary

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/metachris/flashbotsrpc"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) CreateOrAddToFlashbotsBundle(ur *web3_client.UniversalRouterExecCmd, bn string) {
	if a.Bundle.Txs == nil {
		maxTime := ur.Deadline.Uint64()
		a.Bundle = artemis_flashbots.MevTxBundle{
			FlashbotsSendBundleRequest: flashbotsrpc.FlashbotsSendBundleRequest{
				Txs:          []string{},
				BlockNumber:  bn,
				MaxTimestamp: &maxTime,
			},
		}
	}
	a.Bundle.AddTxs(a.OrderedTxs...)
	a.trackTxs(a.OrderedTxs...)
	a.OrderedTxs = []*types.Transaction{}
}

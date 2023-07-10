package artemis_trading_auxiliary

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/metachris/flashbotsrpc"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type AuxiliaryTradingUtils struct {
	artemis_flashbots.FlashbotsClient
	U          *web3_client.UniswapClient
	OrderedTxs []*types.Transaction
}

func InitAuxiliaryTradingUtils(ctx context.Context, nodeURL, network string, acc accounts.Account) AuxiliaryTradingUtils {
	fba := artemis_flashbots.InitFlashbotsClient(ctx, nodeURL, network, &acc)
	wb3 := web3_client.Web3Client{
		Web3Actions: fba.Web3Actions,
	}
	un := web3_client.InitUniswapClient(ctx, wb3)
	return AuxiliaryTradingUtils{
		U:               &un,
		FlashbotsClient: fba,
	}
}

func (a *AuxiliaryTradingUtils) AddTx(tx *types.Transaction) {
	if a.OrderedTxs == nil {
		a.OrderedTxs = []*types.Transaction{}
	}
	a.OrderedTxs = append(a.OrderedTxs, tx)
	a.trackTx(tx)
}

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

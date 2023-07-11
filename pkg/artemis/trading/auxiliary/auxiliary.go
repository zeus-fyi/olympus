package artemis_trading_auxiliary

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_flashbots "github.com/zeus-fyi/olympus/pkg/artemis/trading/flashbots"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type AuxiliaryTradingUtils struct {
	artemis_flashbots.FlashbotsClient
	U          *web3_client.UniswapClient
	Bundle     artemis_flashbots.MevTxBundle
	MevTxGroup MevTxGroup
}

type MevTxGroup struct {
	EventID    int
	OrderedTxs []*types.Transaction
	MevTxs     []artemis_eth_txs.EthTx
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

func (a *AuxiliaryTradingUtils) GetDeadline() *big.Int {
	deadline := int(time.Now().Add(60 * time.Second).Unix())
	sigDeadline := artemis_eth_units.NewBigInt(deadline)
	return sigDeadline
}

func (a *AuxiliaryTradingUtils) addTx(tx *types.Transaction) {
	if a.MevTxGroup.OrderedTxs == nil {
		a.MevTxGroup.OrderedTxs = []*types.Transaction{}
	}
	a.MevTxGroup.OrderedTxs = append(a.MevTxGroup.OrderedTxs, tx)
}

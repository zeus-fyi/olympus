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
	OrderedTxs []TxWithMetadata
	MevTxs     []artemis_eth_txs.EthTx
}

type TxWithMetadata struct {
	TradeType string
	Tx        *types.Transaction
}

func (m *MevTxGroup) GetRawOrderedTxs() []*types.Transaction {
	txSlice := make([]*types.Transaction, len(m.OrderedTxs))
	for i, tx := range m.OrderedTxs {
		txSlice[i] = tx.Tx
	}
	return txSlice
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

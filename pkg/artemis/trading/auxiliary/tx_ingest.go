package artemis_trading_auxiliary

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
)

/*
type TxWithMetadata struct {
	TradeType string
	SignedTx        *types.Transaction
}
*/

// AddTxToBundleGroup adjusts tx for bundle specific gas and adds to bundle group
func (a *AuxiliaryTradingUtils) AddTxToBundleGroup(ctx context.Context, tx *types.Transaction) (artemis_eth_txs.EthTx, error) {
	if a.MevTxGroup.OrderedTxs == nil {
		a.MevTxGroup.OrderedTxs = []TxWithMetadata{}
	}
	txWithMetadata := a.AddTxMetaData(tx)
	mevTx, err := a.packageTxForBundle(ctx, txWithMetadata)
	if err != nil {
		log.Err(err).Msg("error packaging regular tx")
		return artemis_eth_txs.EthTx{}, err
	}
	a.MevTxGroup.MevTxs = append(a.MevTxGroup.MevTxs, mevTx)
	return mevTx, err
}

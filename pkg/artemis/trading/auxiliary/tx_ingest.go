package artemis_trading_auxiliary

import (
	"context"

	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
)

// AddTxToBundleGroup adjusts tx for bundle specific gas and adds to bundle group
func AddTxToBundleGroup(ctx context.Context, txWithMetadata TxWithMetadata, bundle *MevTxGroup) (*MevTxGroup, error) {
	if bundle == nil {
		bundle = &MevTxGroup{
			OrderedTxs: []TxWithMetadata{},
			MevTxs:     []artemis_eth_txs.EthTx{},
		}
	}
	if bundle.OrderedTxs == nil {
		bundle.OrderedTxs = []TxWithMetadata{}
	}
	if bundle.MevTxs == nil {
		bundle.MevTxs = []artemis_eth_txs.EthTx{}
	}

	from, err := web3_actions.GetSender(txWithMetadata.Tx)
	if err != nil {
		log.Err(err).Msg("error getting sender")
		return nil, err
	}
	mevTx, err := packageTxForBundle(ctx, from.String(), txWithMetadata)
	if err != nil {
		log.Err(err).Msg("error packaging regular tx")
		return nil, err
	}
	if txWithMetadata.Permit2Tx.Owner != "" {
		mevTx.Permit2Tx.Permit2Tx = txWithMetadata.Permit2Tx
	}
	bundle.MevTxs = append(bundle.MevTxs, mevTx)
	bundle.OrderedTxs = append(bundle.OrderedTxs, txWithMetadata)
	return bundle, nil
}

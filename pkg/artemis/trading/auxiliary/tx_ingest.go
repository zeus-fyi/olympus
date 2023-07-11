package artemis_trading_auxiliary

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

/*
type TxWithMetadata struct {
	TradeType string
	SignedTx        *types.Transaction
}
*/

// AddTxToBundleGroup adjusts tx for bundle specific gas and adds to bundle group
func (a *AuxiliaryTradingUtils) AddTxToBundleGroup(ctx context.Context, tx *types.Transaction) error {
	if a.MevTxGroup.OrderedTxs == nil {
		a.MevTxGroup.OrderedTxs = []TxWithMetadata{}
	}
	txWithMetadata := a.AddTxMetaData(tx)
	mevTx, err := a.packageTxForBundle(ctx, txWithMetadata)
	if err != nil {
		log.Err(err).Msg("error packaging regular tx")
		return err
	}
	a.MevTxGroup.MevTxs = append(a.MevTxGroup.MevTxs, mevTx)
	return err
}

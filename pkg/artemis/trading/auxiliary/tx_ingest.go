package artemis_trading_auxiliary

import (
	"context"

	"github.com/rs/zerolog/log"
)

/*
type TxWithMetadata struct {
	TradeType string
	SignedTx        *types.Transaction
}
*/

func (a *AuxiliaryTradingUtils) AddTxToBundleGroup(ctx context.Context, txWithMetadata TxWithMetadata) error {
	if a.MevTxGroup.OrderedTxs == nil {
		a.MevTxGroup.OrderedTxs = []TxWithMetadata{}
	}
	signedTx := txWithMetadata.Tx
	mevTx, err := a.packageRegularTx(ctx, signedTx, 0)
	if err != nil {
		log.Err(err).Msg("error packaging regular tx")
		return err
	}
	a.MevTxGroup.MevTxs = append(a.MevTxGroup.MevTxs, mevTx)
	a.MevTxGroup.OrderedTxs = append(a.MevTxGroup.OrderedTxs, txWithMetadata)
	return err
}

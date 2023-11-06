package artemis_eth_rxs

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type EthRx struct {
}

func getBundleTxReceiptQ() string {
	var que = `
        INSERT INTO eth_tx_receipts (
            tx_hash, 
            event_id, 
            status, 
            gas_used, 
            effective_gas_price, 
            cumulative_gas_used, 
            block_hash, 
            block_number,
            transaction_index
        ) 
        VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        )
        ON CONFLICT (tx_hash) 
        DO UPDATE SET 
            status = EXCLUDED.status,
            gas_used = EXCLUDED.gas_used,
            effective_gas_price = EXCLUDED.effective_gas_price,
            cumulative_gas_used = EXCLUDED.cumulative_gas_used,
            block_hash = EXCLUDED.block_hash,
            block_number = EXCLUDED.block_number,
            transaction_index = EXCLUDED.transaction_index
        `
	return que
}

func InsertTxReceipt(ctx context.Context, ethRx artemis_autogen_bases.EthTxReceipts) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = getBundleTxReceiptQ()
	_, err := apps.Pg.Exec(ctx, q.RawQuery, ethRx.TxHash, ethRx.EventID, ethRx.Status, ethRx.GasUsed, ethRx.EffectiveGasPrice, ethRx.CumulativeGasUsed, ethRx.BlockHash, ethRx.BlockNumber, ethRx.TransactionIndex)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertTxReceipt"))
}

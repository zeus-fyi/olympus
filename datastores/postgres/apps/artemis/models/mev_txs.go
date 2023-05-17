package artemis_validator_service_groups_models

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

var ts chronos.Chronos

func InsertMempoolTx(ctx context.Context, tx artemis_autogen_bases.EthMempoolMevTx) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO eth_mempool_mev_tx(protocol_network_id, tx, tx_id, tx_hash, nonce, "from", "to", block_number, tx_flow_prediction)
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				  ON CONFLICT (tx_hash) DO NOTHING;`
	tx.TxID = ts.UnixTimeStampNow()
	_, err := apps.Pg.Exec(ctx, q.RawQuery, tx.ProtocolNetworkID, tx.Tx, tx.TxID, tx.TxHash, tx.Nonce, tx.From, tx.To, tx.BlockNumber, tx.TxFlowPrediction)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertMempoolTx"))
}

func SelectMempoolTxAtBlockNumber(ctx context.Context, protocolID, blockNumber int) (artemis_autogen_bases.EthMempoolMevTxSlice, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT protocol_network_id, tx, tx_id, tx_hash, nonce, "from", "to", block_number, tx_flow_prediction
				  FROM eth_mempool_mev_tx
				  WHERE protocol_network_id = $1 AND block_number = $2
				  `
	log.Debug().Interface("SelectMempoolTxAtBlockNumber", q.LogHeader(ModelName))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, protocolID, blockNumber)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return nil, err
	}
	mempoolTxs := artemis_autogen_bases.EthMempoolMevTxSlice{}
	defer rows.Close()
	for rows.Next() {
		mempoolTx := artemis_autogen_bases.EthMempoolMevTx{}
		rowErr := rows.Scan(
			&mempoolTx.ProtocolNetworkID, &mempoolTx.Tx, &mempoolTx.TxID, &mempoolTx.TxHash, &mempoolTx.Nonce, &mempoolTx.From, &mempoolTx.To, &mempoolTx.BlockNumber, &mempoolTx.TxFlowPrediction,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return nil, rowErr
		}
		mempoolTxs = append(mempoolTxs, mempoolTx)
	}
	return mempoolTxs, misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

func SelectMempoolTxAtMaxBlockNumber(ctx context.Context, protocolID int) (artemis_autogen_bases.EthMempoolMevTxSlice, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_max_block_number AS (
					SELECT MAX(block_number) AS max_block_number
					FROM eth_mempool_mev_tx	
					WHERE protocol_network_id = $1
				)
				  SELECT protocol_network_id, tx, tx_id, tx_hash, nonce, "from", "to", block_number, tx_flow_prediction
				  FROM eth_mempool_mev_tx
				  WHERE protocol_network_id = $1 AND block_number = (SELECT max_block_number FROM cte_max_block_number)
				  `
	log.Debug().Interface("SelectMempoolTxAtBlockNumber", q.LogHeader(ModelName))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, protocolID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ModelName)); returnErr != nil {
		return nil, err
	}
	mempoolTxs := artemis_autogen_bases.EthMempoolMevTxSlice{}
	defer rows.Close()
	for rows.Next() {
		mempoolTx := artemis_autogen_bases.EthMempoolMevTx{}
		rowErr := rows.Scan(
			&mempoolTx.ProtocolNetworkID, &mempoolTx.Tx, &mempoolTx.TxID, &mempoolTx.TxHash, &mempoolTx.Nonce, &mempoolTx.From, &mempoolTx.To, &mempoolTx.BlockNumber, &mempoolTx.TxFlowPrediction,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return nil, rowErr
		}
		mempoolTxs = append(mempoolTxs, mempoolTx)
	}
	return mempoolTxs, misc.ReturnIfErr(err, q.LogHeader(ModelName))
}

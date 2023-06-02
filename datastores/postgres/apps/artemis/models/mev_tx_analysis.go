package artemis_validator_service_groups_models

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertEthMevTxAnalysis(ctx context.Context, txHistory artemis_autogen_bases.EthMevTxAnalysis) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO eth_mev_tx_analysis(gas_used_wei, metadata, tx_hash, trade_method, end_reason, amount_in,
                                amount_out_addr, expected_profit_amount_out, rx_block_number, amount_in_addr,
                                actual_profit_amount_out)
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				  ON CONFLICT (tx_hash) DO NOTHING;`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, txHistory.GetRowValues("InsertEthMevTxAnalysis")...)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertEthMevTxAnalysis"))
}

func SelectEthMevTxAnalysis(ctx context.Context) (artemis_autogen_bases.EthMevTxAnalysisSlice, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  SELECT gas_used_wei, metadata, tx_hash, trade_method, end_reason, amount_in, amount_out_addr,
				         expected_profit_amount_out, rx_block_number, amount_in_addr,
				         actual_profit_amount_out
				  FROM eth_mev_tx_analysis
				  WHERE rx_block_number > 0
				  `
	txAnalysisSlice := artemis_autogen_bases.EthMevTxAnalysisSlice{}
	log.Debug().Interface("SelectEthMevTxAnalysis", q.LogHeader("SelectEthMevTxAnalysis"))
	rows, err := apps.Pg.Query(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectEthMevTxAnalysis")); returnErr != nil {
		return txAnalysisSlice, err
	}
	defer rows.Close()
	for rows.Next() {
		e := artemis_autogen_bases.EthMevTxAnalysis{}
		rowErr := rows.Scan(
			&e.GasUsedWei, &e.Metadata, &e.TxHash, &e.TradeMethod, &e.EndReason, &e.AmountIn, &e.AmountOutAddr, &e.ExpectedProfitAmountOut, &e.RxBlockNumber, &e.AmountInAddr, &e.ActualProfitAmountOut,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return nil, rowErr
		}
		txAnalysisSlice = append(txAnalysisSlice, e)
	}
	return txAnalysisSlice, misc.ReturnIfErr(err, q.LogHeader("SelectEthMevTxAnalysis"))
}

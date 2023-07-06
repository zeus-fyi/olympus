package artemis_mev_models

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const ModelName = "MevAnalysis"

type TradeMethodStats struct {
	TradeMethod string `json:"trade_method"`
	Count       int    `json:"count"`
}

func SelectTradeMethodStatsBySuccess(ctx context.Context) ([]TradeMethodStats, error) {
	return SelectTradeMethodStatsByEndReason(ctx, "success")
}

func SelectTradeMethodStatsByEndReason(ctx context.Context, endReason string) ([]TradeMethodStats, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `	SELECT trade_method, count(*)
					FROM eth_mev_tx_analysis
					WHERE end_reason = $1
					GROUP BY trade_method
				  `
	var txAnalysisSlice []TradeMethodStats
	log.Debug().Interface("SelectEthMevTxAnalysis", q.LogHeader("SelectEthMevTxAnalysis"))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, endReason)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectEthMevTxAnalysis")); returnErr != nil {
		return txAnalysisSlice, err
	}
	defer rows.Close()
	for rows.Next() {
		tms := TradeMethodStats{}
		rowErr := rows.Scan(
			&tms.TradeMethod, &tms.Count,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("TradeMethodStats"))
			return nil, rowErr
		}
		txAnalysisSlice = append(txAnalysisSlice, tms)
	}
	return txAnalysisSlice, nil
}

func InsertEthMevTxAnalysis(ctx context.Context, txHistory artemis_autogen_bases.EthMevTxAnalysis) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO eth_mev_tx_analysis(gas_used_wei, metadata, tx_hash, trade_method, end_reason, amount_in,
                                amount_out_addr, expected_profit_amount_out, rx_block_number, amount_in_addr,
                                actual_profit_amount_out, pair_address)
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
				  ON CONFLICT (tx_hash) DO UPDATE SET
							  gas_used_wei = EXCLUDED.gas_used_wei,
							  metadata = EXCLUDED.metadata,
							  trade_method = EXCLUDED.trade_method,
							  end_reason = EXCLUDED.end_reason,
							  amount_in = EXCLUDED.amount_in,
							  amount_out_addr = EXCLUDED.amount_out_addr,
							  expected_profit_amount_out = EXCLUDED.expected_profit_amount_out,
							  rx_block_number = EXCLUDED.rx_block_number,
							  amount_in_addr = EXCLUDED.amount_in_addr,
				      		  pair_address = EXCLUDED.pair_address,
							  actual_profit_amount_out = EXCLUDED.actual_profit_amount_out;`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, txHistory.GetRowValues("InsertEthMevTxAnalysis")...)
	if err == pgx.ErrNoRows {
		err = nil
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertEthMevTxAnalysis"))
}

type HistoricalAnalysis struct {
	artemis_autogen_bases.EthMempoolMevTx
	artemis_autogen_bases.EthMevTxAnalysis
}

func SelectEthMevTxAnalysisByTxHash(ctx context.Context, txHash string) ([]HistoricalAnalysis, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  SELECT gas_used_wei, metadata, ta.tx_hash, trade_method, end_reason, amount_in, amount_out_addr,
				         expected_profit_amount_out, rx_block_number, amount_in_addr,
				         actual_profit_amount_out, mem.block_number, mem.tx_flow_prediction, mem.nonce, mem.from
				  FROM eth_mev_tx_analysis ta
				  INNER JOIN eth_mempool_mev_tx mem ON mem.tx_hash = ta.tx_hash
				  WHERE ta.tx_hash = $1
				  LIMIT 1000
				  `
	var txAnalysisSlice []HistoricalAnalysis
	log.Debug().Interface("SelectEthMevTxAnalysis", q.LogHeader("SelectEthMevTxAnalysis"))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, txHash)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectEthMevTxAnalysis")); returnErr != nil {
		return txAnalysisSlice, err
	}
	defer rows.Close()
	for rows.Next() {
		e := artemis_autogen_bases.EthMevTxAnalysis{}
		mem := artemis_autogen_bases.EthMempoolMevTx{}
		rowErr := rows.Scan(
			&e.GasUsedWei, &e.Metadata, &e.TxHash, &e.TradeMethod, &e.EndReason, &e.AmountIn, &e.AmountOutAddr, &e.ExpectedProfitAmountOut, &e.RxBlockNumber, &e.AmountInAddr, &e.ActualProfitAmountOut,
			&mem.BlockNumber, &mem.TxFlowPrediction, &mem.Nonce, &mem.From,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return nil, rowErr
		}
		txAnalysisSlice = append(txAnalysisSlice, HistoricalAnalysis{
			EthMempoolMevTx:  mem,
			EthMevTxAnalysis: e,
		})
	}
	return txAnalysisSlice, misc.ReturnIfErr(err, q.LogHeader("SelectEthMevTxAnalysis"))
}

func SelectEthMevTxAnalysis(ctx context.Context) (artemis_autogen_bases.EthMevTxAnalysisSlice, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  SELECT gas_used_wei, metadata, tx_hash, trade_method, end_reason, amount_in, amount_out_addr,
				         expected_profit_amount_out, rx_block_number, amount_in_addr,
				         actual_profit_amount_out
				  FROM eth_mev_tx_analysis
				  WHERE rx_block_number > 0
				  ORDER BY rx_block_number DESC
				  LIMIT 1000
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

package artemis_mev_models

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func SelectEthMevTxAnalysisByPair(ctx context.Context) ([]HistoricalAnalysis, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  SELECT amount_in_addr, amount_out_addr
				  FROM eth_mev_tx_analysis ta
				  WHERE ta.amount_in_addr != ta.amount_out_addr
				  GROUP BY ta.amount_in_addr, ta.amount_out_addr
				  LIMIT 1000
				  `
	var txAnalysisSlice []HistoricalAnalysis
	log.Debug().Interface("SelectEthMevTxAnalysisByPair", q.LogHeader("SelectEthMevTxAnalysisByPair"))
	rows, err := apps.Pg.Query(ctx, q.RawQuery)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectEthMevTxAnalysisByPair")); returnErr != nil {
		return txAnalysisSlice, err
	}
	defer rows.Close()
	for rows.Next() {
		e := artemis_autogen_bases.EthMevTxAnalysis{}
		mem := artemis_autogen_bases.EthMempoolMevTx{}
		rowErr := rows.Scan(&e.AmountInAddr, &e.AmountOutAddr)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return nil, rowErr
		}
		txAnalysisSlice = append(txAnalysisSlice, HistoricalAnalysis{
			EthMempoolMevTx:  mem,
			EthMevTxAnalysis: e,
		})
	}
	return txAnalysisSlice, misc.ReturnIfErr(err, q.LogHeader("SelectEthMevTxAnalysisByPair"))
}

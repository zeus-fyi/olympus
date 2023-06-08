package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type EthMevTxAnalysis struct {
	GasUsedWei              string `db:"gas_used_wei" json:"gasUsedWei"`
	Metadata                string `db:"metadata" json:"metadata"`
	TxHash                  string `db:"tx_hash" json:"txHash"`
	TradeMethod             string `db:"trade_method" json:"tradeMethod"`
	EndReason               string `db:"end_reason" json:"endReason"`
	AmountIn                string `db:"amount_in" json:"amountIn"`
	AmountOutAddr           string `db:"amount_out_addr" json:"amountOutAddr"`
	ExpectedProfitAmountOut string `db:"expected_profit_amount_out" json:"expectedProfitAmountOut"`
	RxBlockNumber           int    `db:"rx_block_number" json:"rxBlockNumber"`
	AmountInAddr            string `db:"amount_in_addr" json:"amountInAddr"`
	ActualProfitAmountOut   string `db:"actual_profit_amount_out" json:"actualProfitAmountOut"`
}
type EthMevTxAnalysisSlice []EthMevTxAnalysis

func (e *EthMevTxAnalysis) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.GasUsedWei, e.Metadata, e.TxHash, e.TradeMethod, e.EndReason, e.AmountIn, e.AmountOutAddr, e.ExpectedProfitAmountOut, e.RxBlockNumber, e.AmountInAddr, e.ActualProfitAmountOut}
	}
	return pgValues
}
func (e *EthMevTxAnalysis) GetTableColumns() (columnValues []string) {
	columnValues = []string{"gas_used_wei", "metadata", "tx_hash", "trade_method", "end_reason", "amount_in", "amount_out_addr", "expected_profit_amount_out", "rx_block_number", "amount_in_addr", "actual_profit_amount_out"}
	return columnValues
}
func (e *EthMevTxAnalysis) GetTableName() (tableName string) {
	tableName = "eth_mev_tx_analysis"
	return tableName
}

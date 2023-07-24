package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type EthTxReceipts struct {
	Status            string `db:"status" json:"status"`
	GasUsed           int    `db:"gas_used" json:"gasUsed"`
	CumulativeGasUsed int    `db:"cumulative_gas_used" json:"cumulativeGasUsed"`
	BlockHash         string `db:"block_hash" json:"blockHash"`
	TransactionIndex  int    `db:"transaction_index" json:"transactionIndex"`
	TxHash            string `db:"tx_hash" json:"txHash"`
	EventID           int    `db:"event_id" json:"eventID"`
	EffectiveGasPrice int    `db:"effective_gas_price" json:"effectiveGasPrice"`
	BlockNumber       int    `db:"block_number" json:"blockNumber"`
}
type EthTxReceiptsSlice []EthTxReceipts

func (e *EthTxReceipts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.Status, e.GasUsed, e.CumulativeGasUsed, e.BlockHash, e.TransactionIndex, e.TxHash, e.EventID, e.EffectiveGasPrice, e.BlockNumber}
	}
	return pgValues
}
func (e *EthTxReceipts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"status", "gas_used", "cumulative_gas_used", "block_hash", "transaction_index", "tx_hash", "event_id", "effective_gas_price", "block_number"}
	return columnValues
}
func (e *EthTxReceipts) GetTableName() (tableName string) {
	tableName = "eth_tx_receipts"
	return tableName
}

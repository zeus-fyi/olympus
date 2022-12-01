package artemis_autogen_bases

import (
	"database/sql"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type EthEcdsaRxs struct {
	RxID                          int            `db:"rx_id" json:"rxID"`
	TxID                          int            `db:"tx_id" json:"txID"`
	GasUsedCumulativeGwei         sql.NullInt64  `db:"gas_used_cumulative_gwei" json:"gasUsedCumulativeGwei"`
	GasUsedCumulativeGweiDecimals sql.NullInt64  `db:"gas_used_cumulative_gwei_decimals" json:"gasUsedCumulativeGweiDecimals"`
	GasUsedGwei                   sql.NullInt64  `db:"gas_used_gwei" json:"gasUsedGwei"`
	GasUsedGweiDecimals           sql.NullInt64  `db:"gas_used_gwei_decimals" json:"gasUsedGweiDecimals"`
	BlockNumber                   sql.NullInt64  `db:"block_number" json:"blockNumber"`
	BlockTimestamp                sql.NullTime   `db:"block_timestamp" json:"blockTimestamp"`
	TxIndex                       sql.NullInt64  `db:"tx_index" json:"txIndex"`
	ContractAddress               sql.NullString `db:"contract_address" json:"contractAddress"`
	Status                        sql.NullString `db:"status" json:"status"`
}
type EthEcdsaRxsSlice []EthEcdsaRxs

func (e *EthEcdsaRxs) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.GasUsedGweiDecimals, e.GasUsedCumulativeGwei, e.BlockNumber, e.RxID, e.GasUsedGwei, e.GasUsedCumulativeGweiDecimals, e.BlockTimestamp, e.TxIndex, e.ContractAddress, e.Status, e.TxID}
	}
	return pgValues
}
func (e *EthEcdsaRxs) GetTableColumns() (columnValues []string) {
	columnValues = []string{"gas_used_gwei_decimals", "gas_used_cumulative_gwei", "block_number", "rx_id", "gas_used_gwei", "gas_used_cumulative_gwei_decimals", "block_timestamp", "tx_index", "contract_address", "status", "tx_id"}
	return columnValues
}
func (e *EthEcdsaRxs) GetTableName() (tableName string) {
	tableName = "eth_ecdsa_rxs"
	return tableName
}

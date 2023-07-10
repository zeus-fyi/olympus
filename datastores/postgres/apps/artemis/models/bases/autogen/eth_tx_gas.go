package artemis_autogen_bases

import (
	"database/sql"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type EthTxGas struct {
	TxHash    string        `db:"tx_hash" json:"txHash"`
	GasPrice  sql.NullInt64 `db:"gasPrice" json:"gasPrice"`
	GasLimit  sql.NullInt64 `db:"gasLimit" json:"gasLimit"`
	GasTipCap sql.NullInt64 `db:"gasTipCap" json:"gasTipCap"`
	GasFeeCap sql.NullInt64 `db:"gasFeeCap" json:"gasFeeCap"`
}
type EthTxGasSlice []EthTxGas

func (e *EthTxGas) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.TxHash, e.GasPrice, e.GasLimit, e.GasTipCap, e.GasFeeCap}
	}
	return pgValues
}
func (e *EthTxGas) GetTableColumns() (columnValues []string) {
	columnValues = []string{"tx_hash", "gasPrice", "gasLimit", "gasTipCap", "gasFeeCap"}
	return columnValues
}
func (e *EthTxGas) GetTableName() (tableName string) {
	tableName = "eth_tx_gas"
	return tableName
}

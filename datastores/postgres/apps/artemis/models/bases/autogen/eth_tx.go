package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type EthTx struct {
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
	TxHash            string `db:"tx_hash" json:"txHash"`
	Nonce             int    `db:"nonce" json:"nonce"`
	From              string `db:"from" json:"from"`
	Type              string `db:"type" json:"type"`
	EventID           int    `db:"event_id" json:"eventID"`
}
type EthTxSlice []EthTx

func (e *EthTx) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.ProtocolNetworkID, e.TxHash, e.Nonce, e.From, e.Type, e.EventID}
	}
	return pgValues
}
func (e *EthTx) GetTableColumns() (columnValues []string) {
	columnValues = []string{"protocol_network_id", "tx_hash", "nonce", "from", "type", "event_id"}
	return columnValues
}
func (e *EthTx) GetTableName() (tableName string) {
	tableName = "eth_tx"
	return tableName
}

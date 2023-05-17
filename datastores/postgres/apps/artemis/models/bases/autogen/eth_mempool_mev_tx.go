package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type EthMempoolMevTx struct {
	TxID              int    `db:"tx_id" json:"txID"`
	To                string `db:"to" json:"to"`
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
	TxFlowPrediction  string `db:"tx_flow_prediction" json:"txFlowPrediction"`
	TxHash            string `db:"tx_hash" json:"txHash"`
	Nonce             int    `db:"nonce" json:"nonce"`
	From              string `db:"from" json:"from"`
	BlockNumber       int    `db:"block_number" json:"blockNumber"`
	Tx                string `db:"tx" json:"tx"`
}
type EthMempoolMevTxSlice []EthMempoolMevTx

func (e *EthMempoolMevTx) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.TxID, e.To, e.ProtocolNetworkID, e.TxFlowPrediction, e.TxHash, e.Nonce, e.From, e.BlockNumber, e.Tx}
	}
	return pgValues
}
func (e *EthMempoolMevTx) GetTableColumns() (columnValues []string) {
	columnValues = []string{"tx_id", "to", "protocol_network_id", "tx_flow_prediction", "tx_hash", "nonce", "from", "block_number", "tx"}
	return columnValues
}
func (e *EthMempoolMevTx) GetTableName() (tableName string) {
	tableName = "eth_mempool_mev_tx"
	return tableName
}

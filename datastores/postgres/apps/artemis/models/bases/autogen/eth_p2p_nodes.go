package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type EthP2PNodes struct {
	ID                int    `db:"id" json:"id"`
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
	Nodes             string `db:"nodes" json:"nodes"`
}
type EthP2pNodesSlice []EthP2PNodes

func (e *EthP2PNodes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.ID, e.ProtocolNetworkID, e.Nodes}
	}
	return pgValues
}
func (e *EthP2PNodes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"id", "protocol_network_id", "nodes"}
	return columnValues
}
func (e *EthP2PNodes) GetTableName() (tableName string) {
	tableName = "eth_p2p_nodes"
	return tableName
}

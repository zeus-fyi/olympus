package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type EthMevAddressFilter struct {
	Address           string `db:"address" json:"address"`
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
}
type EthMevAddressFilterSlice []EthMevAddressFilter

func (e *EthMevAddressFilter) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.Address, e.ProtocolNetworkID}
	}
	return pgValues
}
func (e *EthMevAddressFilter) GetTableColumns() (columnValues []string) {
	columnValues = []string{"address", "protocol_network_id"}
	return columnValues
}
func (e *EthMevAddressFilter) GetTableName() (tableName string) {
	tableName = "eth_mev_address_filter"
	return tableName
}

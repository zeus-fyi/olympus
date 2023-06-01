package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Erc20TokenInfo struct {
	Address           string `db:"address" json:"address"`
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
	BalanceOfSlotNum  int    `db:"balanceOfSlotNum" json:"balanceOfSlotNum"`
}
type Erc20TokenInfoSlice []Erc20TokenInfo

func (e *Erc20TokenInfo) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.Address, e.ProtocolNetworkID, e.BalanceOfSlotNum}
	}
	return pgValues
}
func (e *Erc20TokenInfo) GetTableColumns() (columnValues []string) {
	columnValues = []string{"address", "protocol_network_id", "balanceOfSlotNum"}
	return columnValues
}
func (e *Erc20TokenInfo) GetTableName() (tableName string) {
	tableName = "erc20_token_info"
	return tableName
}

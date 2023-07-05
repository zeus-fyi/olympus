package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type UniswapPairInfo struct {
	TradingEnabled    bool   `db:"trading_enabled" json:"tradingEnabled"`
	Address           string `db:"address" json:"address"`
	FactoryAddress    string `db:"factory_address" json:"factoryAddress"`
	Fee               int    `db:"fee" json:"fee"`
	Version           string `db:"version" json:"version"`
	Token0            string `db:"token0" json:"token0"`
	Token1            string `db:"token1" json:"token1"`
	ProtocolNetworkID int    `db:"protocol_network_id" json:"protocolNetworkID"`
}
type UniswapPairInfoSlice []UniswapPairInfo

func (u *UniswapPairInfo) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{u.TradingEnabled, u.Address, u.FactoryAddress, u.Fee, u.Version, u.Token0, u.Token1, u.ProtocolNetworkID}
	}
	return pgValues
}
func (u *UniswapPairInfo) GetTableColumns() (columnValues []string) {
	columnValues = []string{"trading_enabled", "address", "factory_address", "fee", "version", "token0", "token1", "protocol_network_id"}
	return columnValues
}
func (u *UniswapPairInfo) GetTableName() (tableName string) {
	tableName = "uniswap_pair_info"
	return tableName
}

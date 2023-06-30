package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type Erc20TokenInfo struct {
	Address               string  `db:"address" json:"address"`
	ProtocolNetworkID     int     `db:"protocol_network_id" json:"protocolNetworkID"`
	BalanceOfSlotNum      int     `db:"balanceOfSlotNum" json:"balanceOfSlotNum"`
	Name                  string  `db:"name" json:"name"`
	Symbol                string  `db:"symbol" json:"symbol"`
	Decimals              int     `db:"decimals" json:"decimals"`
	TransferTaxPercentage float64 `db:"transfer_tax_percentage" json:"transferTaxPercentage"`
	TradingEnabled        bool    `db:"trading_enabled" json:"tradingEnabled"`
}
type Erc20TokenInfoSlice []Erc20TokenInfo

func (e *Erc20TokenInfo) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.Address, e.ProtocolNetworkID, e.BalanceOfSlotNum, e.Name, e.Symbol, e.Decimals, e.TransferTaxPercentage, e.TradingEnabled}
	}
	return pgValues
}
func (e *Erc20TokenInfo) GetTableColumns() (columnValues []string) {
	columnValues = []string{"address", "protocol_network_id", "balanceOfSlotNum", "name", "symbol", "decimals", "transfer_tax_percentage", "trading_enabled"}
	return columnValues
}
func (e *Erc20TokenInfo) GetTableName() (tableName string) {
	tableName = "erc20_token_info"
	return tableName
}

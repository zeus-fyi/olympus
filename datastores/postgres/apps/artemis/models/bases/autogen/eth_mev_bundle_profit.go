package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type EthMevBundleProfit struct {
	BundleHash string `db:"bundle_hash" json:"bundleHash"`
	Revenue    int    `db:"revenue" json:"revenue"`
	Costs      int    `db:"costs" json:"costs"`
	Profit     int    `db:"profit" json:"profit"`
}
type EthMevBundleProfitSlice []EthMevBundleProfit

func (e *EthMevBundleProfit) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.BundleHash, e.Revenue, e.Costs, e.Profit}
	}
	return pgValues
}
func (e *EthMevBundleProfit) GetTableColumns() (columnValues []string) {
	columnValues = []string{"bundle_hash", "revenue", "costs", "profit"}
	return columnValues
}
func (e *EthMevBundleProfit) GetTableName() (tableName string) {
	tableName = "eth_mev_bundle_profit"
	return tableName
}

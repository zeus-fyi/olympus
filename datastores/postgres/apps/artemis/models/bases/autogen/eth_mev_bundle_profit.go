package artemis_autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type EthMevBundleProfit struct {
	BundleHash            string `db:"bundle_hash" json:"bundleHash"`
	Revenue               int    `db:"revenue" json:"revenue"`
	RevenuePrediction     int    `db:"revenue_prediction" json:"revenuePrediction"`
	RevenuePredictionSkew int    `db:"revenue_prediction_skew" json:"revenuePredictionSkew"`
	Costs                 int    `db:"costs" json:"costs"`
	Profit                int    `db:"profit" json:"profit"`
}
type EthMevBundleProfitSlice []EthMevBundleProfit

func (e *EthMevBundleProfit) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{e.BundleHash, e.Revenue, e.Costs, e.Profit, e.RevenuePrediction, e.RevenuePredictionSkew}
	}
	return pgValues
}
func (e *EthMevBundleProfit) GetTableColumns() (columnValues []string) {
	columnValues = []string{"bundle_hash", "revenue", "costs", "profit", "revenue_prediction_skew", "revenue_prediction"}
	return columnValues
}
func (e *EthMevBundleProfit) GetTableName() (tableName string) {
	tableName = "eth_mev_bundle_profit"
	return tableName
}

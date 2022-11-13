package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ChartComponentResources struct {
	ChartComponentResourceID int    `db:"chart_component_resource_id" json:"chartComponentResourceID"`
	ChartComponentKindName   string `db:"chart_component_kind_name" json:"chartComponentKindName"`
	ChartComponentApiVersion string `db:"chart_component_api_version" json:"chartComponentApiVersion"`
}
type ChartComponentResourcesSlice []ChartComponentResources

func (c *ChartComponentResources) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartComponentResourceID, c.ChartComponentKindName, c.ChartComponentApiVersion}
	}
	return pgValues
}
func (c *ChartComponentResources) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_component_resource_id", "chart_component_kind_name", "chart_component_api_version"}
	return columnValues
}
func (c *ChartComponentResources) GetTableName() (tableName string) {
	tableName = "chart_component_resources"
	return tableName
}

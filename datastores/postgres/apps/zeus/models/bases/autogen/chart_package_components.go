package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ChartPackageComponents struct {
	ChartSubcomponentParentClassTypeID int `db:"chart_subcomponent_parent_class_type_id" json:"chartSubcomponentParentClassTypeID"`
	ChartPackageID                     int `db:"chart_package_id" json:"chartPackageID"`
}
type ChartPackageComponentsSlice []ChartPackageComponents

func (c *ChartPackageComponents) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentParentClassTypeID, c.ChartPackageID}
	}
	return pgValues
}
func (c *ChartPackageComponents) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_parent_class_type_id", "chart_package_id"}
	return columnValues
}
func (c *ChartPackageComponents) GetTableName() (tableName string) {
	tableName = "chart_package_components"
	return tableName
}

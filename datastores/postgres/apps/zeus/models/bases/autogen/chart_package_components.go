package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ChartPackageComponents struct {
	ChartPackageID                     int `db:"chart_package_id"`
	ChartSubcomponentParentClassTypeID int `db:"chart_subcomponent_parent_class_type_id"`
}
type ChartPackageComponentsSlice []ChartPackageComponents

func (c *ChartPackageComponents) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartPackageID, c.ChartSubcomponentParentClassTypeID}
	}
	return pgValues
}
func (c *ChartPackageComponents) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_package_id", "chart_subcomponent_parent_class_type_id"}
	return columnValues
}
func (c *ChartPackageComponents) GetTableName() (tableName string) {
	tableName = "chart_package_components"
	return tableName
}

package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ChartSubcomponentParentClassTypes struct {
	ChartPackageID                       int    `db:"chart_package_id" json:"chart_package_id"`
	ChartComponentResourceID             int    `db:"chart_component_resource_id" json:"chart_component_resource_id"`
	ChartSubcomponentParentClassTypeID   int    `db:"chart_subcomponent_parent_class_type_id" json:"chart_subcomponent_parent_class_type_id"`
	ChartSubcomponentParentClassTypeName string `db:"chart_subcomponent_parent_class_type_name" json:"chart_subcomponent_parent_class_type_name"`
}
type ChartSubcomponentParentClassTypesSlice []ChartSubcomponentParentClassTypes

func (c *ChartSubcomponentParentClassTypes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartPackageID, c.ChartComponentResourceID, c.ChartSubcomponentParentClassTypeID, c.ChartSubcomponentParentClassTypeName}
	}
	return pgValues
}
func (c *ChartSubcomponentParentClassTypes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_package_id", "chart_component_resource_id", "chart_subcomponent_parent_class_type_id", "chart_subcomponent_parent_class_type_name"}
	return columnValues
}
func (c *ChartSubcomponentParentClassTypes) GetTableName() (tableName string) {
	tableName = "chart_subcomponent_parent_class_types"
	return tableName
}

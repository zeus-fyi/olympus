package autogen_bases

type ChartSubcomponentParentClassTypes struct {
	ChartSubcomponentParentClassTypeName string `db:"chart_subcomponent_parent_class_type_name"`
	ChartPackageID                       int    `db:"chart_package_id"`
	ChartComponentResourceID             int    `db:"chart_component_resource_id"`
	ChartSubcomponentParentClassTypeID   int    `db:"chart_subcomponent_parent_class_type_id"`
}
type ChartSubcomponentParentClassTypesSlice []ChartSubcomponentParentClassTypes

func (c *ChartSubcomponentParentClassTypes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentParentClassTypeName, c.ChartPackageID, c.ChartComponentResourceID, c.ChartSubcomponentParentClassTypeID}
	}
	return pgValues
}
func (c *ChartSubcomponentParentClassTypes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_parent_class_type_name", "chart_package_id", "chart_component_resource_id", "chart_subcomponent_parent_class_type_id"}
	return columnValues
}
func (c *ChartSubcomponentParentClassTypes) GetTableName() (tableName string) {
	tableName = "chart_subcomponent_parent_class_types"
	return tableName
}

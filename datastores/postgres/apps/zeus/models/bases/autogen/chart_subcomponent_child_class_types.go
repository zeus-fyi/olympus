package autogen_bases

type ChartSubcomponentChildClassTypes struct {
	ChartSubcomponentParentClassTypeID  int    `db:"chart_subcomponent_parent_class_type_id"`
	ChartSubcomponentChildClassTypeID   int    `db:"chart_subcomponent_child_class_type_id"`
	ChartSubcomponentChildClassTypeName string `db:"chart_subcomponent_child_class_type_name"`
}
type ChartSubcomponentChildClassTypesSlice []ChartSubcomponentChildClassTypes

func (c *ChartSubcomponentChildClassTypes) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentParentClassTypeID, c.ChartSubcomponentChildClassTypeID, c.ChartSubcomponentChildClassTypeName}
	}
	return pgValues
}
func (c *ChartSubcomponentChildClassTypes) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_parent_class_type_id", "chart_subcomponent_child_class_type_id", "chart_subcomponent_child_class_type_name"}
	return columnValues
}
func (c *ChartSubcomponentChildClassTypes) GetTableName() (tableName string) {
	tableName = "chart_subcomponent_child_class_types"
	return tableName
}

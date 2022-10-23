package autogen_bases

type ChartSubcomponentsChildValues struct {
	ChartSubcomponentChartPackageTemplateInjection bool   `db:"chart_subcomponent_chart_package_template_injection"`
	ChartSubcomponentKeyName                       string `db:"chart_subcomponent_key_name"`
	ChartSubcomponentValue                         string `db:"chart_subcomponent_value"`
	ChartSubcomponentChildClassTypeID              int    `db:"chart_subcomponent_child_class_type_id"`
}
type ChartSubcomponentsChildValuesSlice []ChartSubcomponentsChildValues

func (c *ChartSubcomponentsChildValues) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentChartPackageTemplateInjection, c.ChartSubcomponentKeyName, c.ChartSubcomponentValue, c.ChartSubcomponentChildClassTypeID}
	}
	return pgValues
}
func (c *ChartSubcomponentsChildValues) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_chart_package_template_injection", "chart_subcomponent_key_name", "chart_subcomponent_value", "chart_subcomponent_child_class_type_id"}
	return columnValues
}
func (c *ChartSubcomponentsChildValues) GetTableName() (tableName string) {
	tableName = "chart_subcomponents_child_values"
	return tableName
}

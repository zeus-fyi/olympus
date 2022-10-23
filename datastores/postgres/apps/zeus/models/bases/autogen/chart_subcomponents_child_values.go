package autogen_bases

type ChartSubcomponentsChildValues struct {
	ChartSubcomponentKeyName                       string `db:"chart_subcomponent_key_name"`
	ChartSubcomponentValue                         string `db:"chart_subcomponent_value"`
	ChartSubcomponentChildClassTypeID              int    `db:"chart_subcomponent_child_class_type_id"`
	ChartSubcomponentChartPackageTemplateInjection bool   `db:"chart_subcomponent_chart_package_template_injection"`
}
type ChartSubcomponentsChildValuesSlice []ChartSubcomponentsChildValues

func (c *ChartSubcomponentsChildValues) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentKeyName, c.ChartSubcomponentValue, c.ChartSubcomponentChildClassTypeID, c.ChartSubcomponentChartPackageTemplateInjection}
	}
	return pgValues
}
func (c *ChartSubcomponentsChildValues) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_key_name", "chart_subcomponent_value", "chart_subcomponent_child_class_type_id", "chart_subcomponent_chart_package_template_injection"}
	return columnValues
}
func (c *ChartSubcomponentsChildValues) GetTableName() (tableName string) {
	tableName = "chart_subcomponents_child_values"
	return tableName
}

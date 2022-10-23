package autogen_bases

type ChartSubcomponentsJsonbChildValues struct {
	ChartSubcomponentChildClassTypeID              int    `db:"chart_subcomponent_child_class_type_id"`
	ChartSubcomponentChartPackageTemplateInjection bool   `db:"chart_subcomponent_chart_package_template_injection"`
	ChartSubcomponentFieldName                     string `db:"chart_subcomponent_field_name"`
	ChartSubcomponentJSONbKeyValues                string `db:"chart_subcomponent_jsonb_key_values"`
}
type ChartSubcomponentsJsonbChildValuesSlice []ChartSubcomponentsJsonbChildValues

func (c *ChartSubcomponentsJsonbChildValues) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentChildClassTypeID, c.ChartSubcomponentChartPackageTemplateInjection, c.ChartSubcomponentFieldName, c.ChartSubcomponentJSONbKeyValues}
	}
	return pgValues
}
func (c *ChartSubcomponentsJsonbChildValues) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_child_class_type_id", "chart_subcomponent_chart_package_template_injection", "chart_subcomponent_field_name", "chart_subcomponent_jsonb_key_values"}
	return columnValues
}
func (c *ChartSubcomponentsJsonbChildValues) GetTableName() (tableName string) {
	tableName = "chart_subcomponents_jsonb_child_values"
	return tableName
}

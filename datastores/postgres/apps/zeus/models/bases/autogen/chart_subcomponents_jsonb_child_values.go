package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ChartSubcomponentsJsonbChildValues struct {
	ChartSubcomponentJSONbKeyValues                string `db:"chart_subcomponent_jsonb_key_values" json:"chart_subcomponent_jsonb_key_values"`
	ChartSubcomponentChildClassTypeID              int    `db:"chart_subcomponent_child_class_type_id" json:"chart_subcomponent_child_class_type_id"`
	ChartSubcomponentChartPackageTemplateInjection bool   `db:"chart_subcomponent_chart_package_template_injection" json:"chart_subcomponent_chart_package_template_injection"`
	ChartSubcomponentFieldName                     string `db:"chart_subcomponent_field_name" json:"chart_subcomponent_field_name"`
}
type ChartSubcomponentsJsonbChildValuesSlice []ChartSubcomponentsJsonbChildValues

func (c *ChartSubcomponentsJsonbChildValues) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentJSONbKeyValues, c.ChartSubcomponentChildClassTypeID, c.ChartSubcomponentChartPackageTemplateInjection, c.ChartSubcomponentFieldName}
	}
	return pgValues
}
func (c *ChartSubcomponentsJsonbChildValues) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_jsonb_key_values", "chart_subcomponent_child_class_type_id", "chart_subcomponent_chart_package_template_injection", "chart_subcomponent_field_name"}
	return columnValues
}
func (c *ChartSubcomponentsJsonbChildValues) GetTableName() (tableName string) {
	tableName = "chart_subcomponents_jsonb_child_values"
	return tableName
}

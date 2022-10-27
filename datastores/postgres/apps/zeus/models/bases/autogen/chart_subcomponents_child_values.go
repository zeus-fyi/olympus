package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ChartSubcomponentsChildValues struct {
	ChartSubcomponentChildClassTypeID              int    `db:"chart_subcomponent_child_class_type_id" json:"chart_subcomponent_child_class_type_id"`
	ChartSubcomponentChartPackageTemplateInjection bool   `db:"chart_subcomponent_chart_package_template_injection" json:"chart_subcomponent_chart_package_template_injection"`
	ChartSubcomponentKeyName                       string `db:"chart_subcomponent_key_name" json:"chart_subcomponent_key_name"`
	ChartSubcomponentValue                         string `db:"chart_subcomponent_value" json:"chart_subcomponent_value"`
}
type ChartSubcomponentsChildValuesSlice []ChartSubcomponentsChildValues

func (c *ChartSubcomponentsChildValues) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentChildClassTypeID, c.ChartSubcomponentChartPackageTemplateInjection, c.ChartSubcomponentKeyName, c.ChartSubcomponentValue}
	}
	return pgValues
}
func (c *ChartSubcomponentsChildValues) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_child_class_type_id", "chart_subcomponent_chart_package_template_injection", "chart_subcomponent_key_name", "chart_subcomponent_value"}
	return columnValues
}
func (c *ChartSubcomponentsChildValues) GetTableName() (tableName string) {
	tableName = "chart_subcomponents_child_values"
	return tableName
}

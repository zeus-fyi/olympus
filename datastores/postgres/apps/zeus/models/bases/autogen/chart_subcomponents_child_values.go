package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ChartSubcomponentsChildValues struct {
	ChartSubcomponentChildValuesID                 int    `db:"chart_subcomponent_child_values_id" json:"chartSubcomponentChildValuesID"`
	ChartSubcomponentChildClassTypeID              int    `db:"chart_subcomponent_child_class_type_id" json:"chartSubcomponentChildClassTypeID"`
	ChartSubcomponentChartPackageTemplateInjection bool   `db:"chart_subcomponent_chart_package_template_injection" json:"chartSubcomponentChartPackageTemplateInjection"`
	ChartSubcomponentKeyName                       string `db:"chart_subcomponent_key_name" json:"chartSubcomponentKeyName"`
	ChartSubcomponentValue                         string `db:"chart_subcomponent_value" json:"chartSubcomponentValue"`
}
type ChartSubcomponentsChildValuesSlice []ChartSubcomponentsChildValues

func (c *ChartSubcomponentsChildValues) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ChartSubcomponentChildValuesID, c.ChartSubcomponentChildClassTypeID, c.ChartSubcomponentChartPackageTemplateInjection, c.ChartSubcomponentKeyName, c.ChartSubcomponentValue}
	}
	return pgValues
}
func (c *ChartSubcomponentsChildValues) GetTableColumns() (columnValues []string) {
	columnValues = []string{"chart_subcomponent_child_values_id", "chart_subcomponent_child_class_type_id", "chart_subcomponent_chart_package_template_injection", "chart_subcomponent_key_name", "chart_subcomponent_value"}
	return columnValues
}
func (c *ChartSubcomponentsChildValues) GetTableName() (tableName string) {
	tableName = "chart_subcomponents_child_values"
	return tableName
}

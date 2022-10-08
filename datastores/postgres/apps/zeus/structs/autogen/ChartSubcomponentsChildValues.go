package autogen_structs

type ChartSubcomponentsChildValues struct {
	ChartSubcomponentChildClassTypeID              int    `db:"chart_subcomponent_child_class_type_id"`
	ChartSubcomponentChartPackageTemplateInjection bool   `db:"chart_subcomponent_chart_package_template_injection"`
	ChartSubcomponentKeyName                       string `db:"chart_subcomponent_key_name"`
	ChartSubcomponentValue                         string `db:"chart_subcomponent_value"`
}

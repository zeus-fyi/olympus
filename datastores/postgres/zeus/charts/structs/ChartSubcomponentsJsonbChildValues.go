package models

type ChartSubcomponentsJsonbChildValues struct {
	ChartSubcomponentChildClassTypeID              int    `db:"chart_subcomponent_child_class_type_id"`
	ChartSubcomponentChartPackageTemplateInjection bool   `db:"chart_subcomponent_chart_package_template_injection"`
	ChartSubcomponentFieldName                     string `db:"chart_subcomponent_field_name"`
	ChartSubcomponentJSONbKeyValues                string `db:"chart_subcomponent_jsonb_key_values"`
}

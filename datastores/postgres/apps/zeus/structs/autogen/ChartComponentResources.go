package autogen_structs

type ChartComponentResources struct {
	ChartComponentResourceID int    `db:"chart_component_resource_id"`
	ChartComponentKindName   string `db:"chart_component_kind_name"`
	ChartComponentApiVersion string `db:"chart_component_api_version"`
}

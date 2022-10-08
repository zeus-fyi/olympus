package autogen_structs

type ChartComponentKinds struct {
	ChartComponentKindID     int    `db:"chart_component_kind_id"`
	ChartComponentKindName   string `db:"chart_component_kind_name"`
	ChartComponentApiVersion string `db:"chart_component_api_version"`
}

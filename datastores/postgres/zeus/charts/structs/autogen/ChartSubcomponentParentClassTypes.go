package autogen_structs

type ChartSubcomponentParentClassTypes struct {
	ChartPackageID                       int    `db:"chart_package_id"`
	ChartComponentKindID                 int    `db:"chart_component_kind_id"`
	ChartSubcomponentParentClassTypeID   int    `db:"chart_subcomponent_parent_class_type_id"`
	ChartSubcomponentParentClassTypeName string `db:"chart_subcomponent_parent_class_type_name"`
}

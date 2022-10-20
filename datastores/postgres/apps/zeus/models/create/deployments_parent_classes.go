package create

import "fmt"

const SelectDeploymentResourceID = "(SELECT chart_component_resource_id FROM chart_component_resources WHERE chart_component_kind_name = 'Deployment' AND chart_component_api_version = 'apps/v1')"

func (d *Deployment) addParentClass(pkgId int, pcName string) string {
	s := fmt.Sprintf(
		`INSERT INTO chart_subcomponent_parent_class_types(chart_package_id, chart_component_resource_id, chart_subcomponent_parent_class_type_name)
				 VALUES (%d, %s, '%s')
				 RETURNING chart_subcomponent_parent_class_type_id`, pkgId, SelectDeploymentResourceID, pcName)
	return s
}

func (d *Deployment) insertDeploymentParentClass(pkgId int) string {
	s := fmt.Sprintf(
		`WITH cte_insert_metadata AS (
					%s
				), cte_insert_spec (
					%s
				)`, d.addParentClass(pkgId, "metadata"), d.addParentClass(pkgId, "spec"),
	)
	return s
}

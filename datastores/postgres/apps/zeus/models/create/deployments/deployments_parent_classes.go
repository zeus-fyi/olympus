package deployments

import "fmt"

const SelectDeploymentResourceID = "(SELECT chart_component_resource_id FROM chart_component_resources WHERE chart_component_kind_name = 'Deployment' AND chart_component_api_version = 'apps/v1')"

func (d *Deployment) addParentClass(pkgId int, pcName string) string {
	s := fmt.Sprintf(
		`INSERT INTO chart_subcomponent_parent_class_types(chart_package_id, chart_component_resource_id, chart_subcomponent_parent_class_type_name)
				 VALUES (%d, %s, '%s')
				 RETURNING chart_subcomponent_parent_class_type_id`, pkgId, SelectDeploymentResourceID, pcName)
	return s
}

func (d *Deployment) insertDeploymentCtes(pkgId int) string {
	s := fmt.Sprintf(
		`WITH cte_insert_metadata AS (
					%s
				), cte_insert_spec AS (
					%s
				), `,
		d.addParentClass(pkgId, "metadata"),
		d.addParentClass(pkgId, "spec"),
	)

	s = d.insertDeploymentMetadataChildren(s, "cte_insert_metadata")
	s = d.insertSpecChildren(s, "cte_insert_spec")

	fakeCte := " cte_term AS ( SELECT 1 ) SELECT true"
	returnExpr := s + fakeCte
	return returnExpr
}

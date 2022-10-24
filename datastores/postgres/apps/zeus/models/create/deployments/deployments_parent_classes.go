package deployments

import (
	"fmt"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

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

func CreateChildClassSingleValueSubCTEs() (sql_query_templates.SubCTE, sql_query_templates.SubCTE) {
	subCTE := sql_query_templates.NewSubInsertCTE("cte_insert_container_ports")
	subCTE.TableName = "chart_subcomponents_child_values"
	subCTE.Fields = []string{"chart_subcomponent_child_class_type_id", "chart_subcomponent_key_name", "chart_subcomponent_value"}
	relationshipsSubCTE := sql_query_templates.NewSubInsertCTE("cte_containers_ports_relationship")
	relationshipsSubCTE.TableName = "chart_subcomponent_child_class_types"
	relationshipsSubCTE.Fields = []string{"chart_subcomponent_parent_class_type_id", "chart_subcomponent_child_class_type_name"}
	return subCTE, relationshipsSubCTE
}

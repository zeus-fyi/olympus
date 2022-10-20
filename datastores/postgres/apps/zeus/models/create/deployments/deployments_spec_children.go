package deployments

import "fmt"

func insertDeploymentSpecChildren() string {
	s := fmt.Sprintf(
		`cte_insert_cct AS (
					INSERT INTO chart_subcomponent_child_class_types(chart_subcomponent_parent_class_type_id, chart_subcomponent_child_class_type_name)
					VALUES ((SELECT chart_subcomponent_parent_class_type_id FROM cte_insert_spec), '%s')
					RETURNING chart_subcomponent_child_class_type_id
	)`, "c")
	return s
}

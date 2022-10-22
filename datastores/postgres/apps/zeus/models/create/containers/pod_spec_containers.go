package containers

import "fmt"

func (p *PodContainersGroup) insertContainerToPodSpec(parentExpression, volID, childClassTypeID string) string {
	valsToInsert := "VALUES "
	valsToInsert += fmt.Sprintf("('%s', %s)", childClassTypeID, volID)
	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO chart_subcomponent_spec_pod_template_containers(chart_subcomponent_child_class_type_id, container_id, is_init_container)
					%s
	),`, "cte_containers_volumes_spec_relationship", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}

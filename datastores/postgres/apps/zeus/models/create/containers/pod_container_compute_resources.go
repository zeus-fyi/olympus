package containers

import (
	"fmt"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

func (p *PodContainersGroup) insertContainerComputeResourcesHeader() string {
	return "INSERT INTO container_compute_resources(compute_resources_id, compute_resources_key_values_jsonb) VALUES "
}

// optional, should skip if not specified/nothing is provided
func (p *PodContainersGroup) getContainerComputeResourcesValues(parentExpression string, cr *autogen_structs.ContainerComputeResources) string {
	parentExpression += fmt.Sprintf("('%d', '%s')", cr.ComputeResourcesID, cr.ComputeResourcesKeyValuesJSONb)
	return parentExpression
}

func (p *PodContainersGroup) insertContainerComputeResourcesRelationship(parentExpression, containerImageID string, envVar autogen_structs.ContainerEnvironmentalVars, cct autogen_structs.ChartSubcomponentChildClassTypes) string {
	valsToInsert := "VALUES "
	valsToInsert += fmt.Sprintf("('%d', (%s), '%d')", cct.ChartSubcomponentChildClassTypeID, selectRelatedContainerIDFromImageID(containerImageID), envVar.EnvID)
	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO containers_environmental_vars(chart_subcomponent_child_class_type_id, container_id, env_id)
					%s
	),`, "cte_container_ports", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}

package containers

import (
	"fmt"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

func (p *PodContainersGroup) insertContainerComputeResourcesHeader() string {
	return "INSERT INTO container_compute_resources(compute_resources_id, compute_resources_key_values_jsonb) VALUES "
}

// optional, should skip if not specified/nothing is provided
func (p *PodContainersGroup) getContainerComputeResourcesValues(parentExpression string, cr *autogen_bases.ContainerComputeResources) string {
	parentExpression += fmt.Sprintf("\n('%d', '%s')", cr.ComputeResourcesID, cr.ComputeResourcesKeyValuesJSONb)
	return parentExpression
}

func (p *PodContainersGroup) insertContainerComputeResourcesRelationshipHeader() string {
	return "INSERT INTO containers_environmental_vars(compute_resources_id, compute_resources_key_values_jsonb) VALUES "
}

func (p *PodContainersGroup) insertContainerComputeResourcesRelationship(parentExpression, containerImageID string, envVar autogen_bases.ContainerEnvironmentalVars, cct autogen_bases.ChartSubcomponentChildClassTypes) string {
	valsToInsert := "VALUES "
	valsToInsert += fmt.Sprintf("\n('%d', (%s), '%d')", cct.ChartSubcomponentChildClassTypeID, selectRelatedContainerIDFromImageID(containerImageID), envVar.EnvID)
	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO containers_environmental_vars(chart_subcomponent_child_class_type_id, container_id, env_id)
					%s
	),`, "cte_container_ports", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}

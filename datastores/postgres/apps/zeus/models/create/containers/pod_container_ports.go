package containers

import (
	"fmt"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

// TODO below
func (p *PodContainersGroup) insertContainerResourceRequest(parentExpression string, rr *autogen_structs.ContainerComputeResources, workloadChildGroupInfo autogen_structs.ChartSubcomponentChildClassTypes) string {
	if rr == nil {
		return parentExpression
	}
	valsToInsert := fmt.Sprintf("('%d', (%s))", rr.ComputeResourcesID, rr.ComputeResourcesKeyValuesJSONb)
	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO container_compute_resources(compute_resources_id, compute_resources_key_values_jsonb)
					%s
	),`, "cte_container_compute_resources", valsToInsert)
	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}

func (p *PodContainersGroup) insertContainerResourceRequestRelationship(parentExpression string, rrID, containerImageID string) string {
	valsToInsert := "VALUES "

	valsToInsert += fmt.Sprintf("('%s', (%s))", rrID, selectRelatedContainerIDFromImageID(containerImageID))

	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO containers_compute_resources(compute_resources_id, container_id)
					%s
	),`, "cte_containers_compute_resources", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}

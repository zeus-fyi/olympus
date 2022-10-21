package containers

import (
	"fmt"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

func (p *PodContainersGroup) insertContainerProbesHeader() string {
	return "INSERT INTO container_probes(probe_id, probe_key_values_jsonb) VALUES "
}

func (p *PodContainersGroup) getContainerProbesValuesForInsert(parentExpression string, rr *autogen_structs.ContainerComputeResources) string {
	if rr == nil {
		return parentExpression
	}
	parentExpression += fmt.Sprintf("('%d', (%s))", rr.ComputeResourcesID, rr.ComputeResourcesKeyValuesJSONb)
	return parentExpression
}

func (p *PodContainersGroup) insertContainerProbesRelationship(parentExpression, containerImageID string, probes autogen_structs.ContainersProbes) string {
	valsToInsert := "VALUES "

	valsToInsert += fmt.Sprintf("('%d', (%s), '%s')", probes.ProbeID, selectRelatedContainerIDFromImageID(containerImageID), probes.ProbeType)

	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO containers_probes(probe_id, container_id, probe_type)
					%s
	),`, "cte_containers_probes", valsToInsert)

	returnExpression := fmt.Sprintf("%s %s", parentExpression, containerInsert)
	return returnExpression
}

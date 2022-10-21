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
	parentExpression += fmt.Sprintf("\n('%d', (%s))", rr.ComputeResourcesID, rr.ComputeResourcesKeyValuesJSONb)
	return parentExpression
}

func (p *PodContainersGroup) insertContainerProbesRelationshipHeader() string {
	return "INSERT INTO container_probes(probe_id, container_id, probe_type) VALUES "
}

func (p *PodContainersGroup) getContainerProbesRelationship(parentExpression, imageID string, probes autogen_structs.ContainersProbes) string {
	valsToInsert := fmt.Sprintf("\n('%d', (%s), '%s')", probes.ProbeID, selectRelatedContainerIDFromImageID(imageID), probes.ProbeType)
	returnExpression := fmt.Sprintf("%s %s", parentExpression, valsToInsert)
	return returnExpression
}

package containers

import (
	"fmt"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

func (p *PodContainersGroup) insertContainerProbesHeader() string {
	return "INSERT INTO container_probes(probe_id, probe_key_values_jsonb) VALUES "
}

func (p *PodContainersGroup) getContainerProbesValuesForInsert(parentExpression string, probe *autogen_bases.ContainerProbes) string {
	if probe == nil {
		return parentExpression
	}
	parentExpression += fmt.Sprintf("\n('%d', (%s))", probe.ProbeID, probe.ProbeKeyValuesJSONb)
	return parentExpression
}

func (p *PodContainersGroup) insertContainerProbesRelationshipHeader() string {
	return "INSERT INTO container_probes(probe_id, container_id, probe_type) VALUES "
}

func (p *PodContainersGroup) getContainerProbesRelationship(parentExpression, imageID string, probes autogen_bases.ContainersProbes) string {
	valsToInsert := fmt.Sprintf("\n('%d', (%s), '%s')", probes.ProbeID, selectRelatedContainerIDFromImageID(imageID), probes.ProbeType)
	returnExpression := fmt.Sprintf("%s %s", parentExpression, valsToInsert)
	return returnExpression
}

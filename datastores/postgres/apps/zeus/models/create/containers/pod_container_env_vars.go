package containers

import (
	"fmt"
)

func (p *PodContainersGroup) insertContainerEnvVarsHeader() string {
	return "INSERT INTO container_environmental_vars(port_id, port_name, container_port, host_port) VALUES "
}

func (p *PodContainersGroup) getInsertContainerEnvVarsValues(parentExpression, containerImageID string) string {
	c, ok := p.Containers[containerImageID]
	if !ok {
		return ""
	}
	for _, ev := range c.Env {
		parentExpression += fmt.Sprintf("('%d', '%s', '%s')", ev.EnvID, ev.Name, ev.Value)
	}

	return parentExpression
}

func (p *PodContainersGroup) insertContainerEnvVarRelationshipHeader() string {
	return "INSERT INTO container_environmental_vars(port_id, port_name, container_port, host_port) VALUES "
}

func (p *PodContainersGroup) getContainerEnvVarRelationshipValues(parentExpression, containerImageID, classTypeID string) string {
	valsToInsert := ""
	c, ok := p.Containers[containerImageID]
	if !ok {
		return valsToInsert
	}
	for _, ev := range c.Env {
		parentExpression += fmt.Sprintf("('%s', (%s), '%d'),", classTypeID, selectRelatedContainerIDFromImageID(containerImageID), ev.EnvID)
	}
	returnExpression := fmt.Sprintf("%s %s", parentExpression, valsToInsert)
	return returnExpression
}

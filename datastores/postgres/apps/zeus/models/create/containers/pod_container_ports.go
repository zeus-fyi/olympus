package containers

import (
	"fmt"
)

func (p *PodContainersGroup) insertContainerPortsHeader() string {
	return "INSERT INTO container_ports(port_id, port_name, container_port, host_port) VALUES "
}

func (p *PodContainersGroup) getContainerPortsValuesForInsert(parentExpression, containerImageID string) string {
	c, ok := p.Containers[containerImageID]
	if !ok {
		return ""
	}
	for _, port := range c.Ports {
		parentExpression += fmt.Sprintf("('%d', '%s', %d, %d)", port.PortID, port.PortName, port.ContainerPort, port.HostPort)
	}

	return parentExpression
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

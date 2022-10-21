package containers

import (
	"fmt"
)

func (p *PodContainersGroup) insertContainerPorts(containerImageID string) string {
	c, ok := p.Containers[containerImageID]
	if !ok {
		return ""
	}
	valsToInsert := "VALUES "
	for _, port := range c.Ports {
		valsToInsert += fmt.Sprintf("('%d', '%s', %d, %d)", port.PortID, port.PortName, port.ContainerPort, port.HostPort)

	}
	containerInsert := fmt.Sprintf(`
				%s AS (
					INSERT INTO container_ports(port_id, port_name, container_port, host_port)
					%s
	),`, "cte_container_ports", valsToInsert)
	returnExpression := fmt.Sprintf("%s", containerInsert)
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

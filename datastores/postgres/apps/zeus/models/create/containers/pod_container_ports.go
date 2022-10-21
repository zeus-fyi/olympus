package containers

import (
	"fmt"
)

func (p *PodContainersGroup) insertContainerPortsHeader() string {
	return "INSERT INTO container_ports(port_id, port_name, container_port, host_port) VALUES "
}

func (p *PodContainersGroup) getContainerPortsValuesForInsert(parentExpression, imageID string, isLastValuesGroup bool) string {
	c, ok := p.Containers[imageID]
	if !ok {
		return ""
	}
	for i, port := range c.Ports {
		parentExpression += fmt.Sprintf("\n('%d', '%s', %d, %d)", port.PortID, port.PortName, port.ContainerPort, port.HostPort)
		if i < len(c.Ports)-1 && !isLastValuesGroup {
			parentExpression += ","
		}
	}

	return parentExpression
}

func (p *PodContainersGroup) insertContainerPortsHeaderRelationshipHeader() string {
	return "INSERT INTO containers_ports(chart_subcomponent_child_class_type_id, container_id, port_id) VALUES "
}

func (p *PodContainersGroup) getContainerPortsHeaderRelationshipValues(parentExpression, imageID, childClassTypeID string) string {
	valsToInsert := ""
	c, ok := p.Containers[imageID]
	if !ok {
		return valsToInsert
	}
	for _, port := range c.Ports {
		valsToInsert += fmt.Sprintf("\n('%s', (%s), '%d')", childClassTypeID, selectRelatedContainerIDFromImageID(imageID), port.PortID)
	}
	returnExpression := fmt.Sprintf("%s %s", parentExpression, valsToInsert)
	return returnExpression
}

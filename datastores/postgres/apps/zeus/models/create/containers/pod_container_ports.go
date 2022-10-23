package containers

import (
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (p *PodContainersGroup) insertContainerPortsHeader() string {
	return "INSERT INTO container_ports(port_id, port_name, container_port, host_port) VALUES "
}

func (p *PodContainersGroup) getContainerPortsValuesForInsert(imageID string, cteSubfield sql_query_templates.SubCTE) {
	c, ok := p.Containers[imageID]
	if !ok {
		return
	}
	for _, port := range c.Ports {
		cteSubfield.AddValues(port.PortID, port.PortName, port.ContainerPort, port.HostPort)
	}
	cteSubfield.AddValues()
	return
}

func (p *PodContainersGroup) insertContainerPortsHeaderRelationshipHeader() string {
	return "INSERT INTO containers_ports(chart_subcomponent_child_class_type_id, container_id, port_id) VALUES "
}

func (p *PodContainersGroup) getContainerPortsHeaderRelationshipValues(childClassTypeID int, imageID string, cteSubfield sql_query_templates.SubCTE) {
	c, ok := p.Containers[imageID]
	if !ok {
		return
	}
	for _, port := range c.Ports {
		cteSubfield.AddValues(childClassTypeID, selectRelatedContainerIDFromImageID(imageID), port.PortID)
	}
	return
}

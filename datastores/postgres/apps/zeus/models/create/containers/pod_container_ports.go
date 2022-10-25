package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func (p *PodTemplateSpec) getContainerPortsValuesForInsert(m map[string]containers.Container, imageID string, cteSubfield *sql_query_templates.SubCTE) {
	c, ok := m[imageID]
	if !ok {
		return
	}
	for _, port := range c.GetPorts() {
		cteSubfield.AddValues(port.PortID, port.PortName, port.ContainerPort, port.HostPort)
	}
	return
}

func (p *PodTemplateSpec) insertContainerPortsHeaderRelationshipHeader() string {
	return "INSERT INTO containers_ports(chart_subcomponent_child_class_type_id, container_id, port_id) VALUES "
}

func (p *PodTemplateSpec) getContainerPortsHeaderRelationshipValues(m map[string]containers.Container, imageID string, cteSubfield *sql_query_templates.SubCTE) {
	c, ok := m[imageID]
	if !ok {
		return
	}
	podSpecChildClassTypeID := p.GetPodSpecChildClassTypeID()
	for _, port := range c.GetPorts() {
		cteSubfield.AddValues(podSpecChildClassTypeID, selectRelatedContainerIDFromImageID(imageID), port.PortID)
	}
	return
}

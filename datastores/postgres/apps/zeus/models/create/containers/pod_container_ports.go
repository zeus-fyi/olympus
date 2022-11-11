package containers

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreatePortsCTEs() (sql_query_templates.SubCTE, sql_query_templates.SubCTE) {
	// env vars
	ts := chronos.Chronos{}
	portsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_insert_container_ports_%d", ts.UnixTimeStampNow()))
	portsSubCTE.TableName = "container_ports"
	portsSubCTE.Columns = []string{"port_id", "port_name", "container_port", "host_port"}
	portsRelationshipsSubCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_containers_ports_relationship_%d", ts.UnixTimeStampNow()))
	portsRelationshipsSubCTE.TableName = "containers_ports"
	portsRelationshipsSubCTE.Columns = []string{"chart_subcomponent_child_class_type_id", "container_id", "port_id"}
	return portsSubCTE, portsRelationshipsSubCTE
}

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

func (p *PodTemplateSpec) getContainerPortsHeaderRelationshipValues(m map[string]containers.Container, imageID string, cteSubfield *sql_query_templates.SubCTE) {
	c, ok := m[imageID]
	if !ok {
		return
	}
	podSpecChildClassTypeID := p.GetPodSpecChildClassTypeID()
	for _, port := range c.GetPorts() {
		cteSubfield.AddValues(podSpecChildClassTypeID, c.GetContainerID(), port.PortID)
	}
	return
}

package autogen_bases

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

type ContainerPorts struct {
	PortName      string `db:"port_name" json:"port_name"`
	ContainerPort int    `db:"container_port" json:"container_port"`
	HostIp        string `db:"host_ip" json:"host_ip"`
	HostPort      int    `db:"host_port" json:"host_port"`
	PortProtocol  string `db:"port_protocol" json:"port_protocol"`
	PortID        int    `db:"port_id" json:"port_id"`
}
type ContainerPortsSlice []ContainerPorts

func (c *ContainerPorts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.PortName, c.ContainerPort, c.HostIp, c.HostPort, c.PortProtocol, c.PortID}
	}
	return pgValues
}
func (c *ContainerPorts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"port_name", "container_port", "host_ip", "host_port", "port_protocol", "port_id"}
	return columnValues
}
func (c *ContainerPorts) GetTableName() (tableName string) {
	tableName = "container_ports"
	return tableName
}

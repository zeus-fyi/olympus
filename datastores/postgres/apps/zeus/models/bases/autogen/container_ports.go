package autogen_bases

type ContainerPorts struct {
	ContainerPort int    `db:"container_port"`
	HostIp        string `db:"host_ip"`
	HostPort      int    `db:"host_port"`
	PortProtocol  string `db:"port_protocol"`
	PortID        int    `db:"port_id"`
	PortName      string `db:"port_name"`
}
type ContainerPortsSlice []ContainerPorts

func (c *ContainerPorts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.ContainerPort, c.HostIp, c.HostPort, c.PortProtocol, c.PortID, c.PortName}
	}
	return pgValues
}
func (c *ContainerPorts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"container_port", "host_ip", "host_port", "port_protocol", "port_id", "port_name"}
	return columnValues
}
func (c *ContainerPorts) GetTableName() (tableName string) {
	tableName = "container_ports"
	return tableName
}

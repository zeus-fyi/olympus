package autogen_bases

type ContainerPorts struct {
	HostPort      int    `db:"host_port"`
	PortProtocol  string `db:"port_protocol"`
	PortID        int    `db:"port_id"`
	PortName      string `db:"port_name"`
	ContainerPort int    `db:"container_port"`
	HostIp        string `db:"host_ip"`
}
type ContainerPortsSlice []ContainerPorts

func (c *ContainerPorts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.HostPort, c.PortProtocol, c.PortID, c.PortName, c.ContainerPort, c.HostIp}
	}
	return pgValues
}
func (c *ContainerPorts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"host_port", "port_protocol", "port_id", "port_name", "container_port", "host_ip"}
	return columnValues
}
func (c *ContainerPorts) GetTableName() (tableName string) {
	tableName = "container_ports"
	return tableName
}

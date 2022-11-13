package autogen_bases

type ContainerPorts struct {
	HostIp        string `db:"host_ip" json:"hostIp"`
	HostPort      int    `db:"host_port" json:"hostPort"`
	PortProtocol  string `db:"port_protocol" json:"portProtocol"`
	PortID        int    `db:"port_id" json:"portID"`
	PortName      string `db:"port_name" json:"portName"`
	ContainerPort int    `db:"container_port" json:"containerPort"`
}
type ContainerPortsSlice []ContainerPorts

func (c *ContainerPorts) GetRowValues(queryName string) apps.RowValues {
	pgValues := apps.RowValues{}
	switch queryName {
	default:
		pgValues = apps.RowValues{c.HostIp, c.HostPort, c.PortProtocol, c.PortID, c.PortName, c.ContainerPort}
	}
	return pgValues
}
func (c *ContainerPorts) GetTableColumns() (columnValues []string) {
	columnValues = []string{"host_ip", "host_port", "port_protocol", "port_id", "port_name", "container_port"}
	return columnValues
}
func (c *ContainerPorts) GetTableName() (tableName string) {
	tableName = "container_ports"
	return tableName
}

package autogen_structs

type ContainerPorts struct {
	PortID        int    `db:"port_id"`
	PortName      string `db:"port_name"`
	ContainerPort int    `db:"container_port"`
	HostIp        string `db:"host_ip"`
	HostPort      int    `db:"host_port"`
	PortProtocol  string `db:"port_protocol"`
}

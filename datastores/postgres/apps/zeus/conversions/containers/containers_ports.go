package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	containers "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	v1 "k8s.io/api/core/v1"
)

func ContainerPortsToDB(cps []v1.ContainerPort) containers.Ports {
	contPortsSlice := make(containers.Ports, len(cps))
	for i, p := range cps {
		port := ContainerPortToDB(p)
		contPortsSlice[i] = port
	}
	return contPortsSlice
}

func ContainerPortToDB(p v1.ContainerPort) autogen_structs.ContainerPorts {
	dbPort := autogen_structs.ContainerPorts{
		PortID:        0,
		PortName:      p.Name,
		ContainerPort: int(p.ContainerPort),
		HostIp:        "",
		HostPort:      int(p.HostPort),
		PortProtocol:  "",
	}
	return dbPort
}

func ConvertContainerPortsToContainerDB(cs v1.Container, dbContainer containers.Container) containers.Container {
	dbContainer.Ports = ContainerPortsToDB(cs.Ports)
	return dbContainer
}

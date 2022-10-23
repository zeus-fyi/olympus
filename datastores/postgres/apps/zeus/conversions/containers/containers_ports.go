package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	v1 "k8s.io/api/core/v1"
)

func ContainerPortsToDB(cps []v1.ContainerPort) autogen_bases.ContainerPortsSlice {
	contPortsSlice := make(autogen_bases.ContainerPortsSlice, len(cps))
	for i, p := range cps {
		port := ContainerPortToDB(p)
		contPortsSlice[i] = port
	}
	return contPortsSlice
}

func ContainerPortToDB(p v1.ContainerPort) autogen_bases.ContainerPorts {
	dbPort := autogen_bases.ContainerPorts{
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

package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers/probes"
	v1 "k8s.io/api/core/v1"
)

type Container struct {
	Metadata        autogen_bases.Containers
	VolumeMounts    autogen_bases.ContainerVolumeMountsSlice
	Ports           autogen_bases.ContainerPortsSlice
	Env             autogen_bases.ContainerEnvironmentalVarsSlice
	Probes          probes.ProbeSlice
	ResourceRequest *autogen_bases.ContainerComputeResources
	IsInitContainer bool

	K8sContainer v1.Container
}

func (c *Container) SetContainerID(id int) {
	c.Metadata.ContainerID = id
}

func (c *Container) GetContainerID() int {
	return c.Metadata.ContainerID
}

func (c *Container) GetPorts() autogen_bases.ContainerPortsSlice {
	return c.Ports
}

func (c *Container) GetEnvVars() autogen_bases.ContainerEnvironmentalVarsSlice {
	return c.Env
}

type Containers []Container

func NewContainer() Container {
	c := Container{
		Metadata:        autogen_bases.Containers{},
		VolumeMounts:    autogen_bases.ContainerVolumeMountsSlice{},
		Ports:           autogen_bases.ContainerPortsSlice{},
		Env:             autogen_bases.ContainerEnvironmentalVarsSlice{},
		Probes:          probes.ProbeSlice{},
		IsInitContainer: false,
		K8sContainer:    v1.Container{},
	}
	return c
}

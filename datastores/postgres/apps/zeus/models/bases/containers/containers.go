package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers/probes"
	v1 "k8s.io/api/core/v1"
)

type Container struct {
	K8sContainer v1.Container
	Metadata     autogen_bases.Containers

	SecurityContext autogen_bases.ContainerSecurityContext
	CmdArgs         autogen_bases.ContainerCommandArgs
	VolumeMounts    autogen_bases.ContainerVolumeMountsSlice
	Ports           autogen_bases.ContainerPortsSlice
	Env             autogen_bases.ContainerEnvironmentalVarsSlice
	Probes          probes.ProbeSlice

	ResourceRequest *autogen_bases.ContainerComputeResources

	DB DbContainers
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

func (c *Container) SetIsInitContainer(b bool) {
	c.Metadata.IsInitContainer = b
}

type Containers []Container

func NewContainer() Container {
	c := Container{
		Metadata: autogen_bases.Containers{
			IsInitContainer:          false,
			ContainerID:              0,
			ContainerName:            "",
			ContainerImageID:         "",
			ContainerVersionTag:      "",
			ContainerPlatformOs:      "",
			ContainerRepository:      "",
			ContainerImagePullPolicy: "",
		},
		VolumeMounts: autogen_bases.ContainerVolumeMountsSlice{},
		Ports:        autogen_bases.ContainerPortsSlice{},
		Env:          autogen_bases.ContainerEnvironmentalVarsSlice{},
		Probes:       probes.ProbeSlice{},
		K8sContainer: v1.Container{},
		DB:           DbContainers{},
	}
	return c
}

func NewInit() Container {
	c := Container{
		Metadata: autogen_bases.Containers{
			IsInitContainer:          true,
			ContainerID:              0,
			ContainerName:            "",
			ContainerImageID:         "",
			ContainerVersionTag:      "",
			ContainerPlatformOs:      "",
			ContainerRepository:      "",
			ContainerImagePullPolicy: "",
		},
		VolumeMounts: autogen_bases.ContainerVolumeMountsSlice{},
		Ports:        autogen_bases.ContainerPortsSlice{},
		Env:          autogen_bases.ContainerEnvironmentalVarsSlice{},
		Probes:       probes.ProbeSlice{},
		K8sContainer: v1.Container{},
		DB:           DbContainers{},
	}
	return c
}

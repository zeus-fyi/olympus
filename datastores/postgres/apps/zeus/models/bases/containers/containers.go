package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
)

type Container struct {
	ClassDefinition autogen_bases.ChartSubcomponentChildClassTypes
	Metadata        autogen_bases.Containers
	VolumeMounts    autogen_bases.ContainerVolumeMountsSlice
	Ports           autogen_bases.ContainerPortsSlice
	Env             autogen_bases.ContainerEnvironmentalVarsSlice
	Probes          autogen_bases.ContainerProbesSlice
	ResourceRequest *autogen_bases.ContainerComputeResources
	IsInitContainer bool
}

func (c *Container) GetPorts() autogen_bases.ContainerPortsSlice {
	return c.Ports
}

func (c *Container) GetEnvVars() autogen_bases.ContainerEnvironmentalVarsSlice {
	return c.Env
}

func (c *Container) GetContainerClassDefinition() autogen_bases.ChartSubcomponentChildClassTypes {
	return c.ClassDefinition
}

type Containers []Container

func NewContainer() Container {
	c := Container{
		ClassDefinition: autogen_bases.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: "",
		},
		Metadata:        autogen_bases.Containers{},
		VolumeMounts:    autogen_bases.ContainerVolumeMountsSlice{},
		Ports:           autogen_bases.ContainerPortsSlice{},
		Env:             autogen_bases.ContainerEnvironmentalVarsSlice{},
		Probes:          autogen_bases.ContainerProbesSlice{},
		IsInitContainer: false,
	}
	return c
}

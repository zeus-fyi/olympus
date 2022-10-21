package containers

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

type Container struct {
	ClassDefinition autogen_structs.ChartSubcomponentChildClassTypes
	Metadata        autogen_structs.Containers
	Ports           Ports
	Env             ContainerEnvVars
	Probes          ContainerProbes
	ResourceRequest *autogen_structs.ContainerComputeResources
}

type Containers []Container

func NewContainer() Container {
	c := Container{
		ClassDefinition: autogen_structs.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: "",
		},
		Metadata: autogen_structs.Containers{},
	}
	return c
}

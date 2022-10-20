package containers

import (
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
)

type Container struct {
	ClassDefinition autogen_structs.ChartSubcomponentChildClassTypes
	Information     autogen_structs.Containers
	Ports           ContainersPorts
	Env             ContainerEnvVars
	Probes          ContainerProbes
}

type Containers []Container

func NewContainer() Container {
	c := Container{
		ClassDefinition: autogen_structs.ChartSubcomponentChildClassTypes{
			ChartSubcomponentParentClassTypeID:  0,
			ChartSubcomponentChildClassTypeID:   0,
			ChartSubcomponentChildClassTypeName: "",
		},
		Information: autogen_structs.Containers{},
	}
	return c
}

package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	v1 "k8s.io/api/core/v1"
)

func ContainerEnvVarsToDB(envs []v1.EnvVar) containers.ContainerEnvVars {
	envVarsSlice := make(containers.ContainerEnvVars, len(envs))
	for i, env := range envs {
		envVar := ContainerEnvVarToDB(env)
		envVarsSlice[i] = envVar
	}
	return envVarsSlice
}

func ContainerEnvVarToDB(env v1.EnvVar) autogen_structs.ContainerEnvironmentalVars {
	dbContainer := autogen_structs.ContainerEnvironmentalVars{
		EnvID: 0,
		Name:  env.Name,
		Value: env.Value,
	}
	return dbContainer
}

func ConvertContainerEnvVarsToDB(cs v1.Container, dbContainer containers.Container) containers.Container {
	dbContainer.Env = ContainerEnvVarsToDB(cs.Env)
	return dbContainer
}

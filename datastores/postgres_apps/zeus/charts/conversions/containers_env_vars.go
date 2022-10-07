package conversions

import (
	v1 "k8s.io/api/core/v1"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/charts/structs/autogen"
)

func ContainerEnvVarsToDB(cs *v1.Container) []autogen_structs.ContainerEnvironmentalVars {
	envVarsSlice := make([]autogen_structs.ContainerEnvironmentalVars, len(cs.Env))
	for i, env := range cs.Env {
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

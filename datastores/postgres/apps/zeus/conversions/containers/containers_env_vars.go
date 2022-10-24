package containers

import (
	"encoding/json"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	v1 "k8s.io/api/core/v1"
)

func ContainerEnvVarsToDB(envs []v1.EnvVar) (autogen_bases.ContainerEnvironmentalVarsSlice, error) {
	envVarsSlice := make(autogen_bases.ContainerEnvironmentalVarsSlice, len(envs))
	for i, env := range envs {
		envVar, err := ContainerEnvVarToDB(env)
		if err != nil {
			return envVarsSlice, err
		}
		envVarsSlice[i] = envVar
	}
	return envVarsSlice, nil
}

func ContainerEnvVarToDB(env v1.EnvVar) (autogen_bases.ContainerEnvironmentalVars, error) {

	dbContainer := autogen_bases.ContainerEnvironmentalVars{
		EnvID: 0,
		Name:  env.Name,
	}

	// selects the value from as a second choice to make ownership of the value more clear
	if len(env.Value) <= 0 {
		bytes, err := json.Marshal(env.ValueFrom)
		if err != nil {
			return dbContainer, err
		}
		dbContainer.Value = string(bytes)
	}

	return dbContainer, nil
}

func ConvertContainerEnvVarsToDB(cs v1.Container, dbContainer containers.Container) (containers.Container, error) {
	env, err := ContainerEnvVarsToDB(cs.Env)
	if err != nil {
		return dbContainer, err
	}
	dbContainer.Env = env
	return dbContainer, nil
}

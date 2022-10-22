package containers

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
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

	// TODO fix
	var jsonValues map[string]interface{}
	if len(env.Value) <= 0 {

		bytes, err := env.ValueFrom.Marshal()
		dev_hacks.Use(err)
		err = json.Unmarshal(bytes, &jsonValues)
	}
	dbContainer := autogen_structs.ContainerEnvironmentalVars{
		EnvID: 0,
		Name:  env.Name,
		Value: jsonValues,
	}
	return dbContainer
}

func ConvertContainerEnvVarsToDB(cs v1.Container, dbContainer containers.Container) containers.Container {
	dbContainer.Env = ContainerEnvVarsToDB(cs.Env)
	return dbContainer
}

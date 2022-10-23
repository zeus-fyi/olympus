package containers

import (
	"encoding/json"

	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
	v1 "k8s.io/api/core/v1"
)

func ContainerEnvVarsToDB(envs []v1.EnvVar) autogen_bases.ContainerEnvironmentalVarsSlice {
	envVarsSlice := make(autogen_bases.ContainerEnvironmentalVarsSlice, len(envs))
	for i, env := range envs {
		envVar := ContainerEnvVarToDB(env)
		envVarsSlice[i] = envVar
	}
	return envVarsSlice
}

func ContainerEnvVarToDB(env v1.EnvVar) autogen_bases.ContainerEnvironmentalVars {

	// TODO fix
	var jsonValues map[string]interface{}
	if len(env.Value) <= 0 {

		bytes, err := env.ValueFrom.Marshal()
		dev_hacks.Use(err)
		err = json.Unmarshal(bytes, &jsonValues)
	}
	dbContainer := autogen_bases.ContainerEnvironmentalVars{
		EnvID: 0,
		Name:  env.Name,
		// TODO fix data type Value: jsonValues,
	}
	return dbContainer
}

func ConvertContainerEnvVarsToDB(cs v1.Container, dbContainer containers.Container) containers.Container {
	dbContainer.Env = ContainerEnvVarsToDB(cs.Env)
	return dbContainer
}

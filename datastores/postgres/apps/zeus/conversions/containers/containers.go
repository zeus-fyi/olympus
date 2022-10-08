package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	v1 "k8s.io/api/core/v1"
)

func ConvertContainersToDB(cs []v1.Container) containers.Containers {
	cl := make([]containers.Container, len(cs))
	for i, c := range cs {
		newContainer := containers.NewContainer()
		newContainer = ConvertContainerInfoToDB(c, newContainer)
		newContainer = ConvertContainerPortsToContainerDB(c, newContainer)
		newContainer = ConvertContainerEnvVarsToDB(c, newContainer)
		newContainer = ConvertContainerProbesToDB(c, newContainer)
		cl[i] = newContainer
	}
	return cl
}

func ConvertContainerInfoToDB(cs v1.Container, dbContainer containers.Container) containers.Container {
	dbContainer.Information = autogen_structs.autogen_structs{
		ContainerName:            cs.Name,
		ContainerImageID:         cs.Image,
		ContainerVersionTag:      "",
		ContainerPlatformOs:      "",
		ContainerRepository:      "",
		ContainerImagePullPolicy: string(cs.ImagePullPolicy),
	}
	return dbContainer
}

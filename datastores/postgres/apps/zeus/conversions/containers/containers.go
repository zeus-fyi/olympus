package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	v1 "k8s.io/api/core/v1"
)

func ConvertContainersToDB(cs []v1.Container) (containers.Containers, error) {
	cl := make([]containers.Container, len(cs))
	for i, c := range cs {
		newContainer := containers.NewContainer()
		newContainer = ConvertContainerInfoToDB(c, newContainer)
		newContainer = ConvertContainerPortsToContainerDB(c, newContainer)
		newContainer = ConvertContainerProbesToDB(c, newContainer)
		newContainer, err := ConvertContainerEnvVarsToDB(c, newContainer)
		if err != nil {
			return cl, err
		}
		cl[i] = newContainer
	}
	return cl, nil
}

func ConvertContainerInfoToDB(cs v1.Container, dbContainer containers.Container) containers.Container {
	dbContainer.Metadata = autogen_bases.Containers{
		ContainerName:            cs.Name,
		ContainerImageID:         cs.Image,
		ContainerVersionTag:      "",
		ContainerPlatformOs:      "",
		ContainerRepository:      "",
		ContainerImagePullPolicy: string(cs.ImagePullPolicy),
	}
	return dbContainer
}

package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/autogen"
	v1 "k8s.io/api/core/v1"
)

func ConvertContainerToDB(cs v1.Container) autogen_structs.Containers {
	dbContainer := autogen_structs.Containers{
		ContainerName:            cs.Name,
		ContainerImageID:         cs.Image,
		ContainerVersionTag:      "",
		ContainerPlatformOs:      "",
		ContainerRepository:      "",
		ContainerImagePullPolicy: string(cs.ImagePullPolicy),
	}
	return dbContainer
}

func ConvertContainersToDB(cs []v1.Container) []autogen_structs.Containers {
	cl := make([]autogen_structs.Containers, len(cs))
	for i, c := range cs {
		cl[i] = ConvertContainerToDB(c)
	}
	return cl
}

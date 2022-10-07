package conversions

import (
	v1 "k8s.io/api/core/v1"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/charts/structs/autogen"
)

func ContainerToDB(cs *v1.Container) autogen_structs.Containers {
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

func ContainersToDB(cs []*v1.Container) []autogen_structs.Containers {
	cl := make([]autogen_structs.Containers, len(cs))
	for i, c := range cs {
		cl[i] = ContainerToDB(c)

	}
	return cl
}

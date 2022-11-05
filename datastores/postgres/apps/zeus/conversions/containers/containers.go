package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	v1 "k8s.io/api/core/v1"
)

func ConvertContainersToDB(cs []v1.Container, isInit bool) (containers.Containers, error) {
	cl := make([]containers.Container, len(cs))
	for i, c := range cs {
		newContainer := containers.NewContainer()
		if isInit {
			newContainer = containers.NewInit()
		}
		newContainer = ConvertContainerCmdArgsToContainerDB(c, newContainer)
		newContainer = ConvertContainerResourcesToContainerDB(c, newContainer)
		newContainer = ConvertContainerInfoToDB(c, newContainer)
		newContainer.VolumeMounts = ContainerVolumesToDB(&c)
		newContainer = ConvertContainerPortsToContainerDB(c, newContainer)
		newContainer, err := ConvertContainerProbesToDB(c, newContainer)
		if err != nil {
			return cl, err
		}
		newContainer, err = ConvertContainerEnvVarsToDB(c, newContainer)
		if err != nil {
			return cl, err
		}
		newContainer, err = ConvertContainerSecurityContextToContainerDB(c, newContainer)
		if err != nil {
			return cl, err
		}
		newContainer.ProcessAndSetAmbiguousContainerFieldStatusAndSubfieldIds()
		newContainer.SetIsInitContainer(isInit)
		cl[i] = newContainer
	}
	return cl, nil
}

func ConvertContainerInfoToDB(cs v1.Container, dbContainer containers.Container) containers.Container {
	var ts chronos.Chronos
	dbContainer.Metadata = autogen_bases.Containers{
		ContainerID:              ts.UnixTimeStampNow(),
		ContainerName:            cs.Name,
		ContainerImageID:         cs.Image,
		ContainerVersionTag:      "",
		ContainerPlatformOs:      "",
		ContainerRepository:      "",
		ContainerImagePullPolicy: string(cs.ImagePullPolicy),
	}
	return dbContainer
}

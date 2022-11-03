package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/containers"
	v1 "k8s.io/api/core/v1"
)

func ConvertContainerResourcesToContainerDB(c v1.Container, dbContainer containers.Container) containers.Container {
	computeResources := autogen_bases.ContainerComputeResources{
		ComputeResourcesCpuRequest:              c.Resources.Requests.Cpu().String(),
		ComputeResourcesCpuLimit:                c.Resources.Limits.Cpu().String(),
		ComputeResourcesRamRequest:              c.Resources.Requests.Memory().String(),
		ComputeResourcesRamLimit:                c.Resources.Limits.Memory().String(),
		ComputeResourcesEphemeralStorageRequest: c.Resources.Requests.StorageEphemeral().String(),
		ComputeResourcesEphemeralStorageLimit:   c.Resources.Limits.StorageEphemeral().String(),
	}
	dbContainer.ResourceRequest = &computeResources
	return dbContainer
}

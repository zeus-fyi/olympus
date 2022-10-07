package conversions

import (
	v1 "k8s.io/api/core/v1"

	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/charts/structs/autogen"
)

func ContainerVolumesToDB(cs *v1.Container) []autogen_structs.ContainerVolumeMounts {
	dbVmSlice := make([]autogen_structs.ContainerVolumeMounts, len(cs.VolumeMounts))
	for i, vm := range cs.VolumeMounts {
		dbVm := ContainerVolumeToDB(vm)
		dbVmSlice[i] = dbVm
	}
	return dbVmSlice
}

func ContainerVolumeToDB(vm v1.VolumeMount) autogen_structs.ContainerVolumeMounts {
	dbContainerVolumeMount := autogen_structs.ContainerVolumeMounts{
		VolumeMountID:   0,
		VolumeMountPath: vm.MountPath,
		VolumeName:      vm.Name,
	}
	return dbContainerVolumeMount
}

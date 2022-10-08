package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	v1 "k8s.io/api/core/v1"
)

func ContainerVolumesToDB(cs *v1.Container) []autogen_structs.autogen_structs {
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

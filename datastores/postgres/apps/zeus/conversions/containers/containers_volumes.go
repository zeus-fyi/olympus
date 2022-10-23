package containers

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	v1 "k8s.io/api/core/v1"
)

func ContainerVolumesToDB(cs *v1.Container) autogen_bases.ContainerVolumeMountsSlice {
	dbVmSlice := make([]autogen_bases.ContainerVolumeMounts, len(cs.VolumeMounts))
	for i, vm := range cs.VolumeMounts {
		dbVm := ContainerVolumeToDB(vm)
		dbVmSlice[i] = dbVm
	}
	return dbVmSlice
}

func ContainerVolumeToDB(vm v1.VolumeMount) autogen_bases.ContainerVolumeMounts {
	dbContainerVolumeMount := autogen_bases.ContainerVolumeMounts{
		VolumeMountID:   0,
		VolumeMountPath: vm.MountPath,
		VolumeName:      vm.Name,
	}
	return dbContainerVolumeMount
}

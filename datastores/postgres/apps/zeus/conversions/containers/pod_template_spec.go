package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/containers"
	v1 "k8s.io/api/core/v1"
)

// ConvertPodTemplateSpecConfigToDB PodTemplateSpecConfigToDB has a dependency on chart_subcomponent_child_class_types and containers
func ConvertPodTemplateSpecConfigToDB(ps *v1.PodSpec) (containers.PodTemplateSpec, error) {
	dbPodSpec := containers.NewPodTemplateSpec()

	dbSpecVolumes, err := common.VolumesToDB(ps.Volumes)
	if err != nil {
		return dbPodSpec, err
	}
	dbPodSpec.Spec.PodTemplateSpecVolumes = dbSpecVolumes

	dbContainers, err := ConvertContainersToDB(ps.Containers)
	if err != nil {
		return dbPodSpec, err
	}
	dbPodSpec.Spec.PodTemplateContainers = dbContainers

	return dbPodSpec, nil
}

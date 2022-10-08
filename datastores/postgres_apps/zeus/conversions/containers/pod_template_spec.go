package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/conversions/common"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/containers"
	v1 "k8s.io/api/core/v1"
)

// ConvertPodTemplateSpecConfigToDB PodTemplateSpecConfigToDB has a dependency on chart_subcomponent_child_class_types and containers
func ConvertPodTemplateSpecConfigToDB(ps *v1.PodSpec) containers.PodTemplateSpec {
	dbPodSpec := containers.NewPodTemplateSpec()

	dbSpecVolumes := common.VolumesToDB(ps.Volumes)
	dbPodSpec.Spec.PodTemplateSpecVolumes = dbSpecVolumes

	dbContainers := ConvertContainersToDB(ps.Containers)
	dbPodSpec.Spec.PodTemplateContainers = dbContainers

	return dbPodSpec
}

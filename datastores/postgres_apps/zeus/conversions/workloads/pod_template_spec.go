package workloads

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/conversions/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
	v1 "k8s.io/api/core/v1"
)

// ConvertPodTemplateSpecConfigToDB PodTemplateSpecConfigToDB has a dependency on chart_subcomponent_child_class_types and containers
func ConvertPodTemplateSpecConfigToDB(ps *v1.PodSpec) common.PodTemplateSpec {
	dbPodSpec := common.NewPodTemplateSpec()

	dbSpecVolumes := conversions.VolumesToDB(ps.Volumes)
	dbPodSpec.Spec.PodTemplateSpecVolumes = dbSpecVolumes

	dbContainers := containers.ConvertContainersToDB(ps.Containers)
	dbPodSpec.Spec.PodTemplateContainers = dbContainers
	return dbPodSpec
}

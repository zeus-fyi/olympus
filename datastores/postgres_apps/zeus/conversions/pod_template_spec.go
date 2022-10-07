package conversions

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/structs/common"
	v1 "k8s.io/api/core/v1"
)

// ConvertPodTemplateSpecConfigToDB PodTemplateSpecConfigToDB has a dependency on chart_subcomponent_child_class_types and containers
func ConvertPodTemplateSpecConfigToDB(ps *v1.PodSpec) common.PodTemplateSpec {
	dbPodSpec := common.NewPodTemplateSpec()

	dbSpecVolumes := VolumesToDB(ps.Volumes)
	dbPodSpec.PodTemplateSpecVolumes = dbSpecVolumes

	dbContainers := ConvertContainersToDB(ps.Containers)
	dbPodSpec.PodTemplateContainers = dbContainers
	return dbPodSpec
}

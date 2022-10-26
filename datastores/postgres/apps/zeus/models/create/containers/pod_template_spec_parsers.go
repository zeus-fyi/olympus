package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	cont_conv "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/containers"
	v1 "k8s.io/api/core/v1"
)

// ConvertPodTemplateSpecConfigToDB PodTemplateSpecConfigToDB has a dependency on chart_subcomponent_child_class_types and containers
func (p *PodTemplateSpec) ConvertPodTemplateSpecConfigToDB(ps *v1.PodSpec) (PodTemplateSpec, error) {
	dbPodSpec := NewPodTemplateSpec()

	dbSpecVolumes, err := common_conversions.VolumesToDB(ps.Volumes)
	if err != nil {
		return dbPodSpec, err
	}
	dbSpecContainers, err := cont_conv.ConvertContainersToDB(ps.Containers)
	if err != nil {
		return dbPodSpec, err
	}
	dbPodSpec.Spec.PodTemplateContainers = dbSpecContainers
	dbPodSpec.Spec.PodTemplateSpecVolumes = dbSpecVolumes
	if err != nil {
		return dbPodSpec, err
	}

	return dbPodSpec, nil
}

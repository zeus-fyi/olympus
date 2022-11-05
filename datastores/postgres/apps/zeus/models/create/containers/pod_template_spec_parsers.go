package containers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	cont_conv "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/containers"
	v1 "k8s.io/api/core/v1"
)

// ConvertPodTemplateSpecConfigToDB PodTemplateSpecConfigToDB has a dependency on chart_subcomponent_child_class_types and containers
func (p *PodTemplateSpec) ConvertPodTemplateSpecConfigToDB(ps *v1.PodSpec) error {
	dbSpecVolumes, err := common_conversions.VolumesToDB(ps.Volumes)
	if err != nil {
		return err
	}
	dbSpecContainers, err := cont_conv.ConvertContainersToDB(ps.Containers, false)
	if err != nil {
		return err
	}
	dbSpecInitContainers, err := cont_conv.ConvertContainersToDB(ps.InitContainers, true)
	if err != nil {
		return err
	}

	dbSpecInitContainers = append(dbSpecInitContainers, dbSpecContainers...)
	p.Spec.PodTemplateContainers = dbSpecInitContainers
	p.Spec.PodTemplateSpecVolumes = dbSpecVolumes
	if err != nil {
		return err
	}

	return nil
}

package statefulset

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/common"
	v1 "k8s.io/api/apps/v1"
)

func ConvertStatefulSetSpecConfigToDB(s *v1.StatefulSet) (StatefulSet, error) {
	dbStatefulSet := NewStatefulSet()
	dbStatefulSet.Metadata.Metadata = common_conversions.CreateMetadataByFields(s.Name, s.Annotations, s.Labels)
	spec, err := ConvertStatefulSetSpec(s.Spec)
	if err != nil {
		return dbStatefulSet, err
	}
	dbStatefulSet.Spec = spec
	return dbStatefulSet, nil
}

func ConvertStatefulSetSpec(s v1.StatefulSetSpec) (Spec, error) {
	spec := Spec{
		SpecWorkload: common.NewSpecWorkload(),
		Template:     containers.NewPodTemplateSpec(),
	}

	podTemplateSpec := s.Template.Spec

	dbPodTemplateSpec, err := spec.Template.ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	if err != nil {
		return spec, err
	}
	spec.Template = dbPodTemplateSpec
	return spec, nil
}

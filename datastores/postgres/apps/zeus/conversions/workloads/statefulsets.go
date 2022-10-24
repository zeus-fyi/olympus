package workloads

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/workloads"
	v1 "k8s.io/api/apps/v1"
)

func ConvertStatefulSetSpecConfigToDB(s *v1.StatefulSet) (workloads.StatefulSet, error) {
	dbStatefulSet := workloads.NewStatefulSet()
	dbStatefulSet.Metadata = common.CreateMetadataByFields(s.Name, s.Annotations, s.Labels)
	spec, err := ConvertStatefulSetSpec(s.Spec)
	if err != nil {
		return dbStatefulSet, err
	}
	dbStatefulSet.Spec = spec
	return dbStatefulSet, nil
}

func ConvertStatefulSetSpec(s v1.StatefulSetSpec) (workloads.StatefulSetSpec, error) {
	statefulSetTemplateSpec := s.Template
	podTemplateSpec := statefulSetTemplateSpec.Spec
	dbStatefulSetSpec := workloads.StatefulSetSpec{
		Replicas: 0,
		Selector: common.ConvertSelector(s.Selector),
	}
	dbPodTemplateSpec, err := containers.ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	if err != nil {
		return dbStatefulSetSpec, err
	}
	dbStatefulSetSpec.Template = dbPodTemplateSpec
	return dbStatefulSetSpec, nil
}

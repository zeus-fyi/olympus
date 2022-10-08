package workloads

import (
	common2 "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/containers"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/workloads"
	v1 "k8s.io/api/apps/v1"
)

func ConvertStatefulSetSpecConfigToDB(s *v1.StatefulSet) workloads.StatefulSet {
	dbStatefulSet := workloads.NewStatefulSet()
	dbStatefulSet.Metadata = common2.CreateMetadataByFields(s.Name, s.Annotations, s.Labels)
	dbStatefulSet.Spec = ConvertStatefulSetSpec(s.Spec)
	return dbStatefulSet
}

func ConvertStatefulSetSpec(s v1.StatefulSetSpec) workloads.StatefulSetSpec {
	statefulSetTemplateSpec := s.Template
	podTemplateSpec := statefulSetTemplateSpec.Spec
	dbPodTemplateSpec := containers.ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	dbStatefulSetSpec := workloads.StatefulSetSpec{
		Replicas: 0,
		Template: dbPodTemplateSpec,
		Selector: common2.ConvertSelector(s.Selector),
	}
	return dbStatefulSetSpec
}

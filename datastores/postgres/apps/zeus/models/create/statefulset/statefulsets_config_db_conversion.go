package statefulset

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/apps/v1"
)

func ConvertStatefulSetSpecConfigToDB(s *v1.StatefulSet) (StatefulSet, error) {
	dbStatefulSet := NewStatefulSet()
	dbStatefulSet.Metadata = common.CreateMetadataByFields(s.Name, s.Annotations, s.Labels)
	spec, err := ConvertStatefulSetSpec(s.Spec)
	if err != nil {
		return dbStatefulSet, err
	}
	dbStatefulSet.Spec = spec
	return dbStatefulSet, nil
}

func ConvertStatefulSetSpec(s v1.StatefulSetSpec) (StatefulSetSpec, error) {
	statefulSetTemplateSpec := s.Template
	podTemplateSpec := statefulSetTemplateSpec.Spec

	dbStatefulSetSpec := StatefulSetSpec{
		Selector: common.ConvertSelector(s.Selector),
	}
	dbStatefulSetSpec.Replicas.ChartSubcomponentValue = string_utils.Convert32BitPtrIntToString(s.Replicas)

	dbPodTemplateSpec, err := dbStatefulSetSpec.Template.ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	if err != nil {
		return dbStatefulSetSpec, err
	}
	dbStatefulSetSpec.Template = dbPodTemplateSpec
	return dbStatefulSetSpec, nil
}

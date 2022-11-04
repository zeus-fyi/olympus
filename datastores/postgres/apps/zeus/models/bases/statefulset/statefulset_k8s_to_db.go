package statefulset

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
)

func (s *StatefulSet) ConvertStatefulSetSpecConfigToDB() error {
	dbStatefulSet := NewStatefulSet()
	dbStatefulSet.Metadata.Metadata = common_conversions.CreateMetadataByFields(s.K8sStatefulSet.Name, s.K8sStatefulSet.Annotations, s.K8sStatefulSet.Labels)
	spec, err := s.ConvertStatefulSetSpec()
	if err != nil {
		return err
	}
	dbStatefulSet.Spec = spec
	return nil
}

func (s *StatefulSet) ConvertStatefulSetSpec() (Spec, error) {
	spec := Spec{
		SpecWorkload: structs.NewSpecWorkload(),
		Template:     containers.NewPodTemplateSpec(),
	}
	s.ConvertK8sStatefulSetUpdateStrategyToDB()
	s.ConvertK8sStatefulPodManagementPolicyToDB()
	s.ConvertK8sStatefulServiceNameToDB()

	podTemplateSpec := s.K8sStatefulSet.Spec.Template.Spec
	dbPodTemplateSpec, err := spec.Template.ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	if err != nil {
		return spec, err
	}
	spec.Template = dbPodTemplateSpec
	return spec, nil
}

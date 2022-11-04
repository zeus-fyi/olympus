package statefulset

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (s *StatefulSet) ConvertK8sStatefulSetToDB() error {
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

	s.Spec.Replicas.ChartSubcomponentValue = string_utils.Convert32BitPtrIntToString(s.K8sStatefulSet.Spec.Replicas)
	spec.Selector = structs.NewSelector()

	m := make(map[string]string)
	if s.K8sStatefulSet.Spec.Selector != nil {
		bytes, err := json.Marshal(s.K8sStatefulSet.Spec.Selector)
		if err != nil {
			return spec, err
		}
		selectorString := string(bytes)
		m["selectorString"] = selectorString
		s.Spec.Selector.MatchLabels.AddValues(m)
	}

	s.ConvertK8sStatefulSetUpdateStrategyToDB()
	s.ConvertK8sStatefulPodManagementPolicyToDB()
	s.ConvertK8sStatefulServiceNameToDB()

	err := s.ConvertK8VolumeClaimTemplatesToDB()
	if err != nil {
		return spec, err
	}
	podTemplateSpec := s.K8sStatefulSet.Spec.Template.Spec
	dbPodTemplateSpec, err := spec.Template.ConvertPodTemplateSpecConfigToDB(&podTemplateSpec)
	if err != nil {
		return spec, err
	}
	spec.Template = dbPodTemplateSpec
	return spec, nil
}

package statefulsets

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/structs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/containers"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (s *StatefulSet) ConvertK8sStatefulSetToDB() error {
	s.Metadata.ChartSubcomponentParentClassTypeName = "StatefulSetParentMetadata"
	s.Metadata.Metadata = common_conversions.ConvertMetadata(s.K8sStatefulSet.ObjectMeta)
	s.Metadata.ChartComponentResourceID = StsChartComponentResourceID
	err := s.ConvertK8sStatefulSetSpecToDB()
	if err != nil {
		return err
	}
	return nil
}

func (s *StatefulSet) ConvertK8sStatefulSetSpecToDB() error {
	spec := Spec{
		SpecWorkload: structs.NewSpecWorkload(),
		Template:     containers.NewPodTemplateSpec(),
	}

	replicaCount := string_utils.Convert32BitPtrIntToString(s.K8sStatefulSet.Spec.Replicas)
	s.Spec.Replicas.ChartSubcomponentValue = replicaCount
	spec.Selector = structs.NewSelector()

	m := make(map[string]string)
	if s.K8sStatefulSet.Spec.Selector != nil {
		bytes, err := json.Marshal(s.K8sStatefulSet.Spec.Selector)
		if err != nil {
			return err
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
		log.Err(err)
		return err
	}
	err = s.Spec.Template.ConvertPodTemplateSpecConfigToDB(&s.K8sStatefulSet.Spec.Template.Spec)
	if err != nil {
		log.Err(err)
		return err
	}
	s.Spec.Template.ParentClass.ChartComponentResourceID = StsChartComponentResourceID
	dbPodTemplateSpecMetadata := s.K8sStatefulSet.Spec.Template.GetObjectMeta()
	s.Spec.Template.Metadata.Metadata = common_conversions.CreateMetadataByFields(dbPodTemplateSpecMetadata.GetName(), dbPodTemplateSpecMetadata.GetAnnotations(), dbPodTemplateSpecMetadata.GetLabels())
	s.Spec.Template.Metadata.ChartComponentResourceID = StsChartComponentResourceID
	return nil
}

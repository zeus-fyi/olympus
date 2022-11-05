package statefulset

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
)

func (s *StatefulSet) ParseDBConfigToK8s(pcSlice common_conversions.ParentChildDB) error {
	dbStatefulSet := NewStatefulSet()
	dbStatefulSet.Metadata.Metadata = common_conversions.CreateMetadataByFields(s.K8sStatefulSet.Name, s.K8sStatefulSet.Annotations, s.K8sStatefulSet.Labels)

	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "Spec":
			err := s.ConvertDBSpecToK8s(pc)
			if err != nil {
				return err
			}
		case "StatefulSetParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&s.K8sStatefulSet.ObjectMeta, pc)
		}
	}
	return nil
}

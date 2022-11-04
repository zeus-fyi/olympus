package statefulset

func (s *StatefulSet) ConvertK8VolumeClaimTemplatesToDB() {
	//pvcs := s.K8sStatefulSet.Spec.VolumeClaimTemplates
	// parent needed?
	s.Spec.StatefulSetUpdateStrategy.ChartSubcomponentChildClassTypeName = "StatefulSetVolumeClaimTemplates"

	//for i, pvc := range pvcs {
	//
	//	pvc.
	//		s.Spec.StatefulSetUpdateStrategy.AddKeyValue("type", string(updateType))
	//}
}

package statefulset

func (s *StatefulSet) ConvertK8VolumeClaimTemplatesToDB() error {
	s.Spec.VolumeClaimTemplates.K8sPersistentVolumeClaimSlice = s.K8sStatefulSet.Spec.VolumeClaimTemplates
	err := s.Spec.VolumeClaimTemplates.ConvertK8VolumeClaimTemplateSliceToDB()
	return err
}

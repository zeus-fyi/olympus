package statefulset

func (s *StatefulSet) ConvertK8sStatefulSetUpdateStrategyToDB() {
	updateType := s.K8sStatefulSet.Spec.UpdateStrategy.Type
	s.Spec.StatefulSetUpdateStrategy.ChartSubcomponentChildClassTypeName = "StatefulSetUpdateStrategy"
	s.Spec.StatefulSetUpdateStrategy.AddKeyValue("type", string(updateType))
}

func (s *StatefulSet) ConvertK8sStatefulPodManagementPolicyToDB() {
	podPolicy := s.K8sStatefulSet.Spec.PodManagementPolicy
	s.Spec.PodManagementPolicy.ChartSubcomponentChildClassTypeName = "StatefulSetPodManagementPolicy"
	s.Spec.PodManagementPolicy.AddKeyValue("podManagementPolicy", string(podPolicy))
}

func (s *StatefulSet) ConvertK8sStatefulServiceNameToDB() {
	svcName := s.K8sStatefulSet.Spec.ServiceName
	s.Spec.ServiceName.ChartSubcomponentChildClassTypeName = "StatefulSetServiceName"
	s.Spec.ServiceName.AddKeyValue("serviceName", svcName)
}

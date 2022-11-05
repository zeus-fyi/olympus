package statefulset

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
)

func (s *StatefulSet) ConvertDBSpecToK8s(pcSlice []common_conversions.PC) error {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		value := pc.ChartSubcomponentValue
		fmt.Println(value)
		fmt.Println(subClassName)

		switch subClassName {
		case "StatefulSetUpdateStrategy":
		case "StatefulSetPodManagementPolicy":
		case "StatefulSetServiceName":
		}
	}
	return nil
}

func (s *StatefulSet) ConvertDBStatefulSetUpdateStrategyToK8s() {

	updateType := s.K8sStatefulSet.Spec.UpdateStrategy.Type
	s.Spec.StatefulSetUpdateStrategy.ChartSubcomponentChildClassTypeName = "StatefulSetUpdateStrategy"
	s.Spec.StatefulSetUpdateStrategy.AddKeyValue("type", string(updateType))
}

func (s *StatefulSet) ConvertDBStatefulPodManagementPolicyToK8s() {
	podPolicy := s.K8sStatefulSet.Spec.PodManagementPolicy
	s.Spec.PodManagementPolicy.ChartSubcomponentChildClassTypeName = "StatefulSetPodManagementPolicy"
	s.Spec.PodManagementPolicy.AddKeyValue("podManagementPolicy", string(podPolicy))
}

func (s *StatefulSet) ConvertDBStatefulServiceNameToK8s() {
	svcName := s.K8sStatefulSet.Spec.ServiceName
	s.Spec.ServiceName.ChartSubcomponentChildClassTypeName = "StatefulSetServiceName"
	s.Spec.ServiceName.AddKeyValue("serviceName", svcName)
}

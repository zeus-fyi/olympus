package statefulset

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/apps/v1"
)

func (s *StatefulSet) ConvertDBSpecToK8s(pcSlice []common_conversions.PC) error {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		value := pc.ChartSubcomponentValue
		switch subClassName {
		case "StatefulSetUpdateStrategy":
			switch value {
			case "RollingUpdate":
				strategy := v1.StatefulSetUpdateStrategyType(value)
				s.K8sStatefulSet.Spec.UpdateStrategy.Type = strategy
			case "OnDelete":
				strategy := v1.StatefulSetUpdateStrategyType(value)
				s.K8sStatefulSet.Spec.UpdateStrategy.Type = strategy
			default:
				log.Info().Msgf("StatefulSetUpdateStrategyType: %s not found", value)
			}
		case "StatefulSetPodManagementPolicy":
			switch value {
			case "OrderedReady":
				podPolicy := v1.PodManagementPolicyType(value)
				s.K8sStatefulSet.Spec.PodManagementPolicy = podPolicy
			case "Parallel":
				podPolicy := v1.PodManagementPolicyType(value)
				s.K8sStatefulSet.Spec.PodManagementPolicy = podPolicy
			default:
				log.Info().Msgf("PodManagementPolicyType: %s not found", value)
			}
		case "StatefulSetServiceName":
			s.K8sStatefulSet.Spec.ServiceName = value
		case "replicas":
			s.K8sStatefulSet.Spec.Replicas = string_utils.ConvertStringTo32BitPtrInt(pc.ChartSubcomponentValue)
		case "selector":
			sl, err := db_to_k8s_conversions.ParseLabelSelectorJsonString(pc.ChartSubcomponentValue)
			s.K8sStatefulSet.Spec.Selector = sl
			if err != nil {
				return err
			}
		}
	}
	return nil
}

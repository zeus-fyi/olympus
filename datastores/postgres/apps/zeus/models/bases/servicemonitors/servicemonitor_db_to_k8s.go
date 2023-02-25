package servicemonitors

import (
	"encoding/json"

	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
)

func (s *ServiceMonitor) ParseDBConfigToK8s(pcSlice common_conversions.ParentChildDB) error {
	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "Spec":
			err := s.ConvertDBSpecToK8s(pc)
			if err != nil {
				return err
			}
		case "ServiceMonitorParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&s.K8sServiceMonitor.ObjectMeta, pc)
		}
	}
	return nil
}

func (s *ServiceMonitor) ConvertDBSpecToK8s(pcSlice []common_conversions.PC) error {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		value := pc.ChartSubcomponentValue
		switch subClassName {
		case "ServiceMonitorSpec":
			s.K8sServiceMonitor.Spec = v1.ServiceMonitorSpec{}
			err := json.Unmarshal([]byte(value), &s.K8sServiceMonitor.Spec)
			if err != nil {
				log.Err(err)
				return err
			}
		}
	}
	return nil
}

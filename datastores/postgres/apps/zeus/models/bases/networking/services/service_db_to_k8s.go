package services

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
	v1 "k8s.io/api/core/v1"
)

func (s *Service) ParsePCGroupMap(pcSlice common_conversions.ParentChildDB) error {
	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "Spec":
			err := s.ConvertSpec(pc)
			if err != nil {
				log.Err(err).Msg("error converting service spec")
				return err
			}
		case "ServiceParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&s.K8sService.ObjectMeta, pc)
		}
	}
	return nil
}

func (s *Service) ConvertSpec(pcSlice []common_conversions.PC) error {
	portMap := make(map[int][]common_conversions.PC)
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		ccTypeID := pc.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID

		keyName := pc.ChartSubcomponentKeyName
		value := pc.ChartSubcomponentValue
		switch subClassName {
		case "type":
			s.K8sService.Spec.Type = v1.ServiceType(value)
		case "selector":
			if len(s.K8sService.Spec.Selector) == 0 {
				s.K8sService.Spec.Selector = make(map[string]string)
			}
			s.K8sService.Spec.Selector[keyName] = value
		case "svcPortname":
			tmp := portMap[ccTypeID]
			portMap[ccTypeID] = append(tmp, pc)
		case "clusterIP":
			s.K8sService.Spec.ClusterIP = value
		}
	}
	err := s.ParseSvcPorts(portMap)
	return err
}

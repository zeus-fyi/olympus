package networking

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
)

func (s *Service) ParsePCGroupMap(pcSlice common_conversions.ParentChildDB) error {
	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "ServiceSpec":

		case "ServiceParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&s.K8sService.ObjectMeta, pc)
		}
	}
	return nil
}

func (s *Service) ConvertSpec(pcSlice []common_conversions.PC) error {
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		switch subClassName {
		case "selectorString":
		}
	}
	return nil
}

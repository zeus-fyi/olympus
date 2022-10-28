package ingresses

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

func (i *Ingress) ParseDBConfigToK8s(pcSlice common_conversions.ParentChildDB) error {
	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "Spec":
			err := i.ConvertSpec(pc)
			if err != nil {
				return err
			}
		case "IngressParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&i.K8sIngress.ObjectMeta, pc)
		}
	}
	return nil
}

func (i *Ingress) ConvertSpec(pcSlice []common_conversions.PC) error {
	ingressRulesMap := make(map[int][]common_conversions.PC)
	ingressTLSMap := make(map[int][]common_conversions.PC)

	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		ccTypeID := pc.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID

		keyName := pc.ChartSubcomponentKeyName
		value := pc.ChartSubcomponentValue
		dev_hacks.Use(ccTypeID, keyName, value)
		switch subClassName {
		case "tls":
			// todo i.K8sIngress.Spec.TLS
			err := i.ConvertDBIngressRuleToK8s(ingressRulesMap)
			if err != nil {
				return err
			}
		case "rules":
			err := i.ConvertDBIngressTLSToK8s(ingressTLSMap)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

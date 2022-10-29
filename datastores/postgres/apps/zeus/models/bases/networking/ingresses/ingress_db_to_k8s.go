package ingresses

import (
	"strings"

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
	ingressRulesMap := make(map[string][]common_conversions.PC)
	ingressTLSMap := make(map[string][]common_conversions.PC)

	namesSlice := []string{}
	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		ccTypeID := pc.ChartSubcomponentChildClassTypes.ChartSubcomponentChildClassTypeID

		namesSlice = append(namesSlice, subClassName)
		keyName := pc.ChartSubcomponentKeyName
		value := pc.ChartSubcomponentValue
		dev_hacks.Use(ccTypeID, keyName, value)

		switch subClassName {
		case "ingressClassName":
			i.K8sIngress.Spec.IngressClassName = &value
		}
		if strings.HasPrefix(subClassName, "tls") {
			tmp := ingressTLSMap[subClassName]
			ingressTLSMap[subClassName] = append(tmp, pc)
		}
		if strings.HasPrefix(subClassName, "rules") {
			tmp := ingressRulesMap[subClassName]
			ingressRulesMap[subClassName] = append(tmp, pc)
		}
	}
	err := i.ConvertDBIngressRuleToK8s(ingressRulesMap)
	if err != nil {
		return err
	}
	err = i.ConvertDBIngressTLSToK8s(ingressTLSMap)
	if err != nil {
		return err
	}
	return nil
}

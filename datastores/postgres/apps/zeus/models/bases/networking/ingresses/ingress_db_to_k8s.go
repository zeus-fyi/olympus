package ingresses

import (
	"strings"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions/db_to_k8s_conversions"
)

func (i *Ingress) ParseDBConfigToK8s(pcSlice common_conversions.ParentChildDB) error {
	for pcGroupName, pc := range pcSlice.PCGroupMap {
		switch pcGroupName {
		case "Spec":
			err := i.ConvertDBSpecToK8s(pc)
			if err != nil {
				return err
			}
		case "IngressParentMetadata":
			db_to_k8s_conversions.ConvertMetadata(&i.K8sIngress.ObjectMeta, pc)
		}
	}
	return nil
}

func (i *Ingress) ConvertDBSpecToK8s(pcSlice []common_conversions.PC) error {
	ingressRulesMap := make(map[string][]common_conversions.PC)
	ingressTLSMap := make(map[string][]common_conversions.PC)

	for _, pc := range pcSlice {
		subClassName := pc.ChartSubcomponentChildClassTypeName
		value := pc.ChartSubcomponentValue

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

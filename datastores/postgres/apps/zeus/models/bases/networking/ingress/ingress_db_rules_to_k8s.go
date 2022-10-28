package ingress

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	v1 "k8s.io/api/networking/v1"
)

func (i *Ingress) ParseToK8sIngressRules(rulesMap map[int][]common_conversions.PC) error {
	for _, rule := range rulesMap {
		k8sIngressRule := v1.IngressRule{}
		for _, ruleComponent := range rule {
			key := ruleComponent.ChartSubcomponentKeyName
			val := ruleComponent.ChartSubcomponentValue
			switch key {
			case "host":
				k8sIngressRule.Host = val
			case "http":

				if k8sIngressRule.IngressRuleValue.HTTP == nil {
					httpRules := v1.HTTPIngressRuleValue{}
					httpRules.Paths = []v1.HTTPIngressPath{}
					k8sIngressRule.IngressRuleValue.HTTP = &httpRules
				}

				// TODO parse

				/* path types
				PathTypeImplementationSpecific = PathType("ImplementationSpecific")
				PathTypePrefix = PathType("Prefix")
				PathTypeExact = PathType("Exact")

				resource if needed
				Resource *v1.TypedLocalObjectReference
				*/

				// TODO parse into this
				ingressBackendService := v1.IngressServiceBackend{
					Name: "",
					Port: v1.ServiceBackendPort{
						Name:   "",
						Number: 0,
					},
				}
				path := v1.HTTPIngressPath{
					Path:     "",
					PathType: nil,
					Backend: v1.IngressBackend{
						Service:  &ingressBackendService,
						Resource: nil,
					},
				}
				paths := k8sIngressRule.IngressRuleValue.HTTP.Paths
				k8sIngressRule.IngressRuleValue.HTTP.Paths = append(paths, path)
			}
		}

	}
	return nil
}

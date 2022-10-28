package ingresses

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	v1 "k8s.io/api/networking/v1"
)

func (i *Ingress) ConvertDBIngressRuleToK8s(rulesMap map[int][]common_conversions.PC) error {
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
				httpPath, err := parseIngressPath(val)
				if err != nil {
					return err
				}
				paths := k8sIngressRule.IngressRuleValue.HTTP.Paths
				k8sIngressRule.IngressRuleValue.HTTP.Paths = append(paths, httpPath)
			}
		}

	}
	return nil
}

func parseIngressPath(ingressPathStr string) (v1.HTTPIngressPath, error) {
	ingressPath := v1.HTTPIngressPath{}
	bytes, berr := getBytes(ingressPathStr)
	if berr != nil {
		return ingressPath, berr
	}
	perr := json.Unmarshal(bytes, &ingressPath)
	if perr != nil {
		return ingressPath, perr
	}

	return ingressPath, nil
}

/* Reference
			PathTypeImplementationSpecific = PathType("ImplementationSpecific")
			PathTypePrefix = PathType("Prefix")
			PathTypeExact = PathType("Exact")

			resource if needed
			Resource *v1.TypedLocalObjectReference

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
*/

func getBytes(structToBytes interface{}) ([]byte, error) {
	bytes, berr := json.Marshal(structToBytes)
	return bytes, berr
}

package ingresses

import "encoding/json"

func (i *Ingress) ConvertK8sIngressRuleToDB() error {
	for _, k8sIngressRule := range i.K8sIngress.Spec.Rules {
		hostName := k8sIngressRule.Host
		httpRules := k8sIngressRule.HTTP
		var paths []string
		if httpRules != nil {
			for _, path := range httpRules.Paths {
				bytes, err := json.Marshal(path)
				if err != nil {
					return err
				}
				paths = append(paths, string(bytes))
			}
			i.Rules.AddIngressRule(hostName, paths)
		}
	}
	return nil
}

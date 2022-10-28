package ingresses

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	v1 "k8s.io/api/networking/v1"
)

func (i *Ingress) ConvertDBIngressTLSToK8s(tlsMap map[int][]common_conversions.PC) error {
	for _, dbTLS := range tlsMap {
		k8sTLS := v1.IngressTLS{
			Hosts:      []string{},
			SecretName: "",
		}
		for _, tlsComponent := range dbTLS {
			val := tlsComponent.ChartSubcomponentValue
			k := tlsComponent.ChartSubcomponentKeyName
			switch k {
			case "secretName":
				k8sTLS.SecretName = val
			case "hosts":
				k8sTLS.Hosts = append(k8sTLS.Hosts, val)
			}
		}
		i.K8sIngress.Spec.TLS = append(i.K8sIngress.Spec.TLS, k8sTLS)
	}
	return nil
}

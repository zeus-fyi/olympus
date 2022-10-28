package ingresses

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/common_conversions"
	v1 "k8s.io/api/networking/v1"
)

func (i *Ingress) ParseToK8sTLS(pcSlice []common_conversions.PC) error {
	// TODO parse into
	tls := v1.IngressTLS{
		Hosts:      nil,
		SecretName: "",
	}
	i.K8sIngress.Spec.TLS = append(i.K8sIngress.Spec.TLS, tls)
	return nil
}

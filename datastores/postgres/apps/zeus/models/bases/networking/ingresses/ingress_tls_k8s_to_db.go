package ingresses

func (i *Ingress) ConvertK8sIngressTLSToDB() {
	for _, k8sTLS := range i.K8sIngress.Spec.TLS {
		secretNameValue := k8sTLS.SecretName
		i.TLS.AddIngressTLS(secretNameValue, k8sTLS.Hosts)
	}
	return
}

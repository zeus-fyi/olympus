package zeus_templates

import (
	"context"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
metadata:
  name: "zeus-client"
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: "nginx"
  tls:
    - secretName: zeus-client-tls
      hosts:
        - host.zeus.fyi
  rules:
    - host: host.zeus.fyi
*/

func GetIngressTemplate(ctx context.Context, name string) *v1.Ingress {
	ingressClassName := "nginx"

	annotations := make(map[string]string)
	annotations["cert-manager.io/cluster-issuer"] = "letsencrypt-prod"

	md := metav1.ObjectMeta{
		Name:        GetIngressName(ctx, name),
		Annotations: annotations,
	}
	return &v1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: md,
		Spec: v1.IngressSpec{
			IngressClassName: &ingressClassName,
			TLS: []v1.IngressTLS{{
				SecretName: GetIngressSecretName(ctx, name),
			}},
		},
	}
}

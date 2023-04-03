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

func GetIngressTemplate(ctx context.Context) *v1.Ingress {
	return &v1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.IngressSpec{},
		Status:     v1.IngressStatus{},
	}
}

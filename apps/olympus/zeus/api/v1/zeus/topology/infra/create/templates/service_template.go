package zeus_templates

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
metadata:
  name: zeus-service
  labels:
    app.kubernetes.io/name: zeus-client
    app.kubernetes.io/instance: zeus-client
    app.kubernetes.io/managed-by: zeus
spec:
  type: ClusterIP
  ports:
  selector:
    app.kubernetes.io/name: zeus-client
    app.kubernetes.io/instance: zeus-client
*/

func GetServiceTemplate(ctx context.Context, name string) *v1.Service {
	labels := GetLabels(ctx, name)
	selectors := GetSelector(ctx, name)
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: v1.ServiceSpec{
			Selector: selectors,
			Type:     v1.ServiceTypeClusterIP,
		},
	}
}

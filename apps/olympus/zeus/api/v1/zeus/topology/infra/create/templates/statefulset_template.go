package zeus_templates

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
metadata:
  name: zeus-client
  labels:
    app.kubernetes.io/name: zeus-client
    app.kubernetes.io/instance: zeus-client
    app.kubernetes.io/managed-by: zeus
  annotations:
    {}
spec:
  podManagementPolicy: OrderedReady
  replicas: 0
  selector:
    matchLabels:
      app.kubernetes.io/name: zeus-client
      app.kubernetes.io/instance: zeus-client
  serviceName: zeus-client
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app.kubernetes.io/name: zeus-client
        app.kubernetes.io/instance: zeus-client
    spec:
      initContainers:
      containers:
      nodeSelector:
        {}
      affinity:
        {}
      tolerations:
        []
      volumes:
  volumeClaimTemplates:
*/

func GetStatefulSetTemplate(ctx context.Context) *v1.StatefulSet {
	return &v1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.StatefulSetSpec{},
	}
}

package zeus_templates

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
metadata:
  name: "zeus-client"
  labels:
    app.kubernetes.io/instance: "zeus-client"
    app.kubernetes.io/name: "zeus-client"
spec:
  replicas: 0
  selector:
    matchLabels:
      app.kubernetes.io/name: "zeus-client"
      app.kubernetes.io/instance: "zeus-client"
  template:
    metadata:
      labels:
        app.kubernetes.io/name: "zeus-client"
*/

func GetDeploymentTemplate(ctx context.Context) *v1.Deployment {
	return &v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       v1.DeploymentSpec{},
	}
}

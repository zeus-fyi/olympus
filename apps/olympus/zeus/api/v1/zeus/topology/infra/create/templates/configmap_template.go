package zeus_templates

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
metadata:
  name: cm-zeus
*/

func GetConfigMapTemplate(ctx context.Context) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Immutable:  nil,
		Data:       nil,
		BinaryData: nil,
	}
}

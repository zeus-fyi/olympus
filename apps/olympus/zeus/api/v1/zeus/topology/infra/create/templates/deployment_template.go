package zeus_templates

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	v1Core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeploymentTemplate(ctx context.Context, name string) *v1.Deployment {
	return &v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   GetDeploymentName(ctx, name),
			Labels: GetLabels(ctx, name),
		},
		Spec: v1.DeploymentSpec{
			Selector: metav1.SetAsLabelSelector(GetSelector(ctx, name)),
			Template: v1Core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: GetLabels(ctx, name),
				},
				Spec: v1Core.PodSpec{},
			},
			Strategy: v1.DeploymentStrategy{},
		},
	}
}

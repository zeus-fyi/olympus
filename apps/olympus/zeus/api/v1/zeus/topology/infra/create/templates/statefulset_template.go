package zeus_templates

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	v1Core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetStatefulSetTemplate(ctx context.Context, name string) *v1.StatefulSet {
	labels := GetLabels(ctx, name)
	selectors := GetSelector(ctx, name)
	return &v1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   GetStatefulSetName(ctx, name),
			Labels: labels,
		},
		Spec: v1.StatefulSetSpec{
			Selector: metav1.SetAsLabelSelector(selectors),
			Template: v1Core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1Core.PodSpec{},
			},
			ServiceName:         GetServiceName(ctx, name),
			PodManagementPolicy: v1.OrderedReadyPodManagement,
			UpdateStrategy: v1.StatefulSetUpdateStrategy{
				Type: v1.RollingUpdateStatefulSetStrategyType,
			},
		},
	}
}

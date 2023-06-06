package zeus_templates

import (
	"context"

	v1Batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetCronJobTemplate(ctx context.Context, name string) *v1Batch.Job {
	return &v1Batch.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CronJob",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: GetCronJobName(ctx, name),
		},
	}
}

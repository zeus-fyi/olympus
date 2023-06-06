package zeus_core

import (
	"context"

	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetCronJobsList(ctx context.Context, kns zeus_common_types.CloudCtxNs) (*v1.CronJobList, error) {
	k.SetContext(kns.Context)
	return k.kc.BatchV1().CronJobs(kns.Namespace).List(ctx, metav1.ListOptions{})
}

func (k *K8Util) CreateCronJob(ctx context.Context, kns zeus_common_types.CloudCtxNs, job *v1.CronJob) (*v1.CronJob, error) {
	k.SetContext(kns.Context)
	job, err := k.kc.BatchV1().CronJobs(kns.Namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (k *K8Util) DeleteCronJob(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string) error {
	k.SetContext(kns.Context)
	err := k.kc.BatchV1().CronJobs(kns.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

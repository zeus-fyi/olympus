package zeus_core

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetJobsList(ctx context.Context, kns zeus_common_types.CloudCtxNs) (*v1.JobList, error) {
	k.SetContext(kns.Context)
	return k.kc.BatchV1().Jobs(kns.Namespace).List(ctx, metav1.ListOptions{})
}

func (k *K8Util) CreateJob(ctx context.Context, kns zeus_common_types.CloudCtxNs, job *v1.Job) (*v1.Job, error) {
	k.SetContext(kns.Context)
	job, err := k.kc.BatchV1().Jobs(kns.Namespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (k *K8Util) DeleteJob(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string) error {
	k.SetContext(kns.Context)
	err := k.kc.BatchV1().Jobs(kns.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		log.Err(err).Msg("DeleteJob")
		return err
	}
	return nil
}

package zeus_core

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetDeploymentList(ctx context.Context, kns zeus_common_types.CloudCtxNs, filter *string_utils.FilterOpts) (*v1.DeploymentList, error) {
	d, err := k.kc.AppsV1().Deployments(kns.Namespace).List(ctx, metav1.ListOptions{})
	return d, err
}

func (k *K8Util) GetDeployment(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Deployment, error) {
	d, err := k.kc.AppsV1().Deployments(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
	return d, err
}

func (k *K8Util) CreateDeployment(ctx context.Context, kns zeus_common_types.CloudCtxNs, d *v1.Deployment, filter *string_utils.FilterOpts) (*v1.Deployment, error) {
	opts := metav1.CreateOptions{}
	d, err := k.kc.AppsV1().Deployments(kns.Namespace).Create(ctx, d, opts)
	return d, err
}

func (k *K8Util) DeleteDeployment(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) error {
	opts := metav1.DeleteOptions{}
	err := k.kc.AppsV1().Deployments(kns.Namespace).Delete(ctx, name, opts)
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) CreateDeploymentIfVersionLabelChangesOrDoesNotExist(ctx context.Context, kns zeus_common_types.CloudCtxNs, nd *v1.Deployment, filter *string_utils.FilterOpts) (*v1.Deployment, error) {
	cd, err := k.GetDeployment(ctx, kns, nd.Name, filter)
	switch {
	case cd != nil && len(cd.Name) > 0:
		switch IsVersionNew(cd.Labels, nd.Labels) {
		case true:
			derr := k.DeleteDeployment(ctx, kns, cd.Name, filter)
			if derr != nil {
				return cd, derr
			}
		case false:
			return cd, nil
		}
	case errors.IsNotFound(err):
		newD, newDErr := k.CreateDeployment(ctx, kns, nd, filter)
		return newD, newDErr
	}
	newD, newDErr := k.CreateDeployment(ctx, kns, nd, filter)
	return newD, newDErr
}

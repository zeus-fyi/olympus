package zeus_core

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetDeployment(ctx context.Context, kns KubeCtxNs, name string) (*v1.Deployment, error) {
	d, err := k.kc.AppsV1().Deployments(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
	return d, err
}

func (k *K8Util) CreateDeployment(ctx context.Context, kns KubeCtxNs, d *v1.Deployment) (*v1.Deployment, error) {
	opts := metav1.CreateOptions{}
	d, err := k.kc.AppsV1().Deployments(kns.Namespace).Create(ctx, d, opts)
	return d, err
}

func (k *K8Util) DeleteDeployment(ctx context.Context, kns KubeCtxNs, name string) error {
	opts := metav1.DeleteOptions{}
	err := k.kc.AppsV1().Deployments(kns.Namespace).Delete(ctx, name, opts)
	return err
}

func (k *K8Util) CreateDeploymentIfVersionLabelChangesOrDoesNotExist(ctx context.Context, kns KubeCtxNs, nd *v1.Deployment) (*v1.Deployment, error) {
	cd, err := k.GetDeployment(ctx, kns, nd.Name)
	switch {
	case cd != nil && len(cd.Name) > 0:
		switch IsVersionNew(cd.Labels, nd.Labels) {
		case true:
			derr := k.DeleteDeployment(ctx, kns, cd.Name)
			if derr != nil {
				return cd, derr
			}
		case false:
			return cd, nil
		}
	case errors.IsNotFound(err):
		newD, newDErr := k.CreateDeployment(ctx, kns, nd)
		return newD, newDErr
	}
	newD, newDErr := k.CreateDeployment(ctx, kns, nd)
	return newD, newDErr
}

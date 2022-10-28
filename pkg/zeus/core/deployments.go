package zeus_core

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetDeployment(ctx context.Context, kubeCtxNs KubeCtxNs, name string) (*v1.Deployment, error) {
	d, err := k.kc.AppsV1().Deployments(kubeCtxNs.Namespace).Get(context.Background(), name, metav1.GetOptions{})
	return d, err
}

func (k *K8Util) CreateDeployment(ctx context.Context, kubeCtxNs KubeCtxNs, d *v1.Deployment) (*v1.Deployment, error) {
	k.SetContext(kubeCtxNs.GetCtxName(kubeCtxNs.Env))
	opts := metav1.CreateOptions{}
	d, err := k.kc.AppsV1().Deployments(kubeCtxNs.Namespace).Create(context.Background(), d, opts)
	return d, err
}

func (k *K8Util) DeleteDeployment(ctx context.Context, kubeCtxNs KubeCtxNs, d *v1.Deployment) (*v1.Deployment, error) {
	k.SetContext(kubeCtxNs.GetCtxName(kubeCtxNs.Env))
	opts := metav1.CreateOptions{}
	d, err := k.kc.AppsV1().Deployments(kubeCtxNs.Namespace).Create(context.Background(), d, opts)
	return d, err
}

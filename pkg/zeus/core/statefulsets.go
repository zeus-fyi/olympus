package zeus_core

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetStatefulSetList(ctx context.Context, kubeCtxNs KubeCtxNs) (*v1.StatefulSetList, error) {
	opts := metav1.ListOptions{}
	ssl, err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).List(ctx, opts)
	return ssl, err
}

func (k *K8Util) GetStatefulSet(ctx context.Context, name string, kubeCtxNs KubeCtxNs) (*v1.StatefulSet, error) {
	k.PrintPath = "stateful_sets"
	k.FileName = name
	opts := metav1.GetOptions{}
	ss, err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).Get(ctx, name, opts)
	if err != nil {
		return ss, err
	}
	_, err = k.K8Printer(ss, kubeCtxNs.Env)
	return ss, err
}

func (k *K8Util) DeleteStatefulSet(ctx context.Context, name string, kubeCtxNs KubeCtxNs) error {
	opts := metav1.DeleteOptions{}
	err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).Delete(ctx, name, opts)
	return err
}

func (k *K8Util) UpdateStatefulSet(ctx context.Context, ss *v1.StatefulSet, kubeCtxNs KubeCtxNs) (*v1.StatefulSet, error) {
	opts := metav1.UpdateOptions{}
	ss, err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).Update(ctx, ss, opts)
	return ss, err
}

func (k *K8Util) CreateStatefulSet(ctx context.Context, ss *v1.StatefulSet, kubeCtxNs KubeCtxNs) (*v1.StatefulSet, error) {
	opts := metav1.CreateOptions{}
	ss, err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).Create(ctx, ss, opts)
	return ss, err
}

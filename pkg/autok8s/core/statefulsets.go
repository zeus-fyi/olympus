package autok8s_core

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetStatefulSetList(kubeCtxNs KubeCtxNs) (*v1.StatefulSetList, error) {
	k.SetContext(kubeCtxNs.GetCtxName(kubeCtxNs.Env))
	opts := metav1.ListOptions{}
	ssl, err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).List(context.Background(), opts)
	return ssl, err
}

func (k *K8Util) GetStatefulSet(name string, kubeCtxNs KubeCtxNs) (*v1.StatefulSet, error) {
	k.PrintPath = "stateful_sets"
	k.FileName = name
	k.SetContext(kubeCtxNs.GetCtxName(kubeCtxNs.Env))
	opts := metav1.GetOptions{}
	ss, err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).Get(context.Background(), name, opts)
	if err != nil {
		return ss, err
	}
	_, err = k.K8Printer(ss, kubeCtxNs.Env)
	return ss, err
}

func (k *K8Util) DeleteStatefulSet(name string, kubeCtxNs KubeCtxNs) error {
	k.SetContext(kubeCtxNs.GetCtxName(kubeCtxNs.Env))
	opts := metav1.DeleteOptions{}
	err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).Delete(context.Background(), name, opts)
	return err
}

func (k *K8Util) UpdateStatefulSet(ss *v1.StatefulSet, kubeCtxNs KubeCtxNs) (*v1.StatefulSet, error) {
	k.SetContext(kubeCtxNs.GetCtxName(kubeCtxNs.Env))
	opts := metav1.UpdateOptions{}
	ss, err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).Update(context.Background(), ss, opts)
	return ss, err
}

func (k *K8Util) CreateStatefulSet(ss *v1.StatefulSet, kubeCtxNs KubeCtxNs) (*v1.StatefulSet, error) {
	k.SetContext(kubeCtxNs.GetCtxName(kubeCtxNs.Env))
	opts := metav1.CreateOptions{}
	ss, err := k.kc.AppsV1().StatefulSets(kubeCtxNs.Namespace).Create(context.Background(), ss, opts)
	return ss, err
}

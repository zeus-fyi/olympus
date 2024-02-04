package zeus_core

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetNamespaces(ctx context.Context, kns zeus_common_types.CloudCtxNs) (*v1.NamespaceList, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
}

func (k *K8Util) CreateNamespace(ctx context.Context, kns zeus_common_types.CloudCtxNs, namespace *v1.Namespace) (*v1.Namespace, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
}

func (k *K8Util) DeleteNamespace(ctx context.Context, kns zeus_common_types.CloudCtxNs) error {
	k.SetContext(kns.Context)
	err := k.kc.CoreV1().Namespaces().Delete(ctx, kns.Namespace, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		log.Err(err).Interface("kns", kns).Msg("DeleteNamespace: error")
		return err
	}
	return err
}

func (k *K8Util) GetNamespace(ctx context.Context, kns zeus_common_types.CloudCtxNs) (*v1.Namespace, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Namespaces().Get(ctx, kns.Namespace, metav1.GetOptions{})
}

func (k *K8Util) CreateNamespaceIfDoesNotExist(ctx context.Context, kns zeus_common_types.CloudCtxNs) (*v1.Namespace, error) {
	k.SetContext(kns.Context)
	ns, err := k.GetNamespace(ctx, kns)
	if errors.IsNotFound(err) {
		ns.Name = kns.Namespace
		ns, err = k.kc.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
		alreadyExists := errors.IsAlreadyExists(err)
		if alreadyExists {
			return ns, nil
		}
		return ns, err
	}
	return ns, err
}

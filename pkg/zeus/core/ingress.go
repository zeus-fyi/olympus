package zeus_core

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetIngressWithKns(ctx context.Context, kns KubeCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Ingress, error) {
	return k.kc.NetworkingV1().Ingresses(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *K8Util) CreateIngressWithKns(ctx context.Context, kns KubeCtxNs, ing *v1.Ingress, filter *string_utils.FilterOpts) (*v1.Ingress, error) {
	return k.kc.NetworkingV1().Ingresses(kns.Namespace).Create(ctx, ing, metav1.CreateOptions{})
}

func (k *K8Util) DeleteIngressWithKns(ctx context.Context, kns KubeCtxNs, name string, filter *string_utils.FilterOpts) error {
	return k.kc.NetworkingV1().Ingresses(kns.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

package zeus_core

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetSecretWithKns(ctx context.Context, kns KubeCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Secret, error) {
	return k.kc.CoreV1().Secrets(kns.Namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (k *K8Util) CreateSecretWithKns(ctx context.Context, kns KubeCtxNs, s *v1.Secret, filter *string_utils.FilterOpts) (*v1.Secret, error) {
	return k.kc.CoreV1().Secrets(kns.Namespace).Create(context.Background(), s, metav1.CreateOptions{})
}

func (k *K8Util) DeleteSecretWithKns(ctx context.Context, kns KubeCtxNs, name string, filter *string_utils.FilterOpts) error {
	return k.kc.CoreV1().Secrets(kns.Namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

func (k *K8Util) CopySecretToAnotherKns(ctx context.Context, kns KubeCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Secret, error) {
	s, err := k.GetSecretWithKns(ctx, kns, name, filter)
	if err != nil {
		return s, err
	}
	s.ResourceVersion = ""
	return k.CreateSecretWithKns(ctx, kns, s, filter)
}

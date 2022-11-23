package zeus_core

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetSecretWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Secret, error) {
	return k.kc.CoreV1().Secrets(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *K8Util) CreateSecretWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, s *v1.Secret, filter *string_utils.FilterOpts) (*v1.Secret, error) {
	return k.kc.CoreV1().Secrets(kns.Namespace).Create(ctx, s, metav1.CreateOptions{})
}

func (k *K8Util) DeleteSecretWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) error {
	err := k.kc.CoreV1().Secrets(kns.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) CopySecretToAnotherKns(ctx context.Context, knsFrom, knsTo zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Secret, error) {
	s, err := k.GetSecretWithKns(ctx, knsFrom, name, filter)
	if err != nil {
		return s, err
	}
	s.ResourceVersion = ""
	s.Namespace = knsTo.Namespace
	return k.CreateSecretWithKns(ctx, knsTo, s, filter)
}

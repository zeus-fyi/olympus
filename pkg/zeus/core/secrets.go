package zeus_core

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetSecretWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Secret, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Secrets(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *K8Util) CreateSecretWithKnsIfDoesNotExist(ctx context.Context, kns zeus_common_types.CloudCtxNs, s *v1.Secret, filter *string_utils.FilterOpts) (*v1.Secret, error) {
	k.SetContext(kns.Context)
	sec, err := k.GetSecretWithKns(ctx, kns, s.Name, nil)
	if errors.IsNotFound(err) {
		return k.CreateSecretWithKns(ctx, kns, s, nil)
	}
	return sec, err
}

func (k *K8Util) CreateSecretWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, s *v1.Secret, filter *string_utils.FilterOpts) (*v1.Secret, error) {
	k.SetContext(kns.Context)
	sec, err := k.kc.CoreV1().Secrets(kns.Namespace).Create(ctx, s, metav1.CreateOptions{})
	alreadyExists := errors.IsAlreadyExists(err)
	if alreadyExists {
		log.Ctx(ctx).Err(err).Msg("Secret already exists, skipping creation")
		return sec, nil
	}
	return sec, err
}

func (k *K8Util) DeleteSecretWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) error {
	k.SetContext(kns.Context)
	err := k.kc.CoreV1().Secrets(kns.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) CopySecretToAnotherKns(ctx context.Context, knsFrom, knsTo zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Secret, error) {
	k.SetContext(knsFrom.Context)
	s, err := k.GetSecretWithKns(ctx, knsFrom, name, filter)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return s, err
	}
	s.ResourceVersion = ""
	s.Namespace = knsTo.Namespace
	k.SetContext(knsTo.Context)
	return k.CreateSecretWithKns(ctx, knsTo, s, filter)
}

package zeus_core

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateSecretWrapper(sec *v1.Secret, kns zeus_common_types.CloudCtxNs, secretName, key, value string) *v1.Secret {
	if sec == nil {
		sec = &v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: kns.Namespace,
			},
			Type: "Opaque",
		}
	}
	if sec.StringData == nil {
		sec.StringData = make(map[string]string)
	}
	sec.StringData[key] = value
	return sec
}

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
		log.Err(err).Msg("Secret already exists, skipping creation")
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

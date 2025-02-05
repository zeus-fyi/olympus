package zeus_core

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetIngressListWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, filter *string_utils.FilterOpts) (*v1.IngressList, error) {
	k.SetContext(kns.Context)
	return k.kc.NetworkingV1().Ingresses(kns.Namespace).List(ctx, metav1.ListOptions{})
}

func (k *K8Util) GetIngressWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Ingress, error) {
	k.SetContext(kns.Context)
	return k.kc.NetworkingV1().Ingresses(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *K8Util) CreateIngressWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, ing *v1.Ingress, filter *string_utils.FilterOpts) (*v1.Ingress, error) {
	k.SetContext(kns.Context)
	ing, err := k.kc.NetworkingV1().Ingresses(kns.Namespace).Create(ctx, ing, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		log.Err(err).Interface("kns", kns).Msg("Ingress already exists, skipping creation")
		return ing, nil
	}
	return ing, err
}

func (k *K8Util) DeleteIngressWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) error {
	k.SetContext(kns.Context)
	err := k.kc.NetworkingV1().Ingresses(kns.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) CreateIngressIfVersionLabelChangesOrDoesNotExist(ctx context.Context, kns zeus_common_types.CloudCtxNs, ning *v1.Ingress, filter *string_utils.FilterOpts) (*v1.Ingress, error) {
	k.SetContext(kns.Context)
	cing, err := k.GetIngressWithKns(ctx, kns, ning.Name, filter)
	switch {
	case cing != nil && len(cing.Name) > 0:
		switch IsVersionNew(cing.Labels, ning.Labels) {
		case true:
			derr := k.DeleteIngressWithKns(ctx, kns, cing.Name, filter)
			if derr != nil {
				return cing, derr
			}
		case false:
			return cing, nil
		}
	case errors.IsNotFound(err):
		newCm, newCmErr := k.CreateIngressWithKns(ctx, kns, ning, filter)
		return newCm, newCmErr
	}
	newCm, newSErr := k.CreateIngressWithKns(ctx, kns, ning, filter)
	return newCm, newSErr
}

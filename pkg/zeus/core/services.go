package zeus_core

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetServiceListWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, filter *string_utils.FilterOpts) (*v1.ServiceList, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Services(kns.Namespace).List(ctx, metav1.ListOptions{})
}

func (k *K8Util) GetServiceWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.Service, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Services(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *K8Util) CreateServiceWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, s *v1.Service, filter *string_utils.FilterOpts) (*v1.Service, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Services(kns.Namespace).Create(ctx, s, metav1.CreateOptions{})
}

func (k *K8Util) DeleteServiceWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) error {
	k.SetContext(kns.Context)
	err := k.kc.CoreV1().Services(kns.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) CreateServiceIfVersionLabelChangesOrDoesNotExist(ctx context.Context, kns zeus_common_types.CloudCtxNs, ns *v1.Service, filter *string_utils.FilterOpts) (*v1.Service, error) {
	k.SetContext(kns.Context)
	cs, err := k.GetServiceWithKns(ctx, kns, ns.Name, filter)
	switch {
	case cs != nil && len(cs.Name) > 0:
		switch IsVersionNew(cs.Labels, ns.Labels) {
		case true:
			derr := k.DeleteServiceWithKns(ctx, kns, cs.Name, filter)
			if derr != nil {
				return cs, derr
			}
		case false:
			return cs, nil
		}
	case errors.IsNotFound(err):
		newS, newSErr := k.CreateServiceWithKns(ctx, kns, ns, filter)
		return newS, newSErr
	}
	newS, newSErr := k.CreateServiceWithKns(ctx, kns, ns, filter)
	return newS, newSErr
}

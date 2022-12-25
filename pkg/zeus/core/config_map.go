package zeus_core

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetConfigMapListWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, filter *string_utils.FilterOpts) (*v1.ConfigMapList, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().ConfigMaps(kns.Namespace).List(ctx, metav1.ListOptions{})
}

func (k *K8Util) GetConfigMapWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) (*v1.ConfigMap, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().ConfigMaps(kns.Namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k *K8Util) CreateConfigMapWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, cm *v1.ConfigMap, filter *string_utils.FilterOpts) (*v1.ConfigMap, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().ConfigMaps(kns.Namespace).Create(ctx, cm, metav1.CreateOptions{})
}

func (k *K8Util) DeleteConfigMapWithKns(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, filter *string_utils.FilterOpts) error {
	k.SetContext(kns.Context)
	err := k.kc.CoreV1().ConfigMaps(kns.Namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}

func (k *K8Util) CreateConfigMapIfVersionLabelChangesOrDoesNotExist(ctx context.Context, kns zeus_common_types.CloudCtxNs, ncm *v1.ConfigMap, filter *string_utils.FilterOpts) (*v1.ConfigMap, error) {
	k.SetContext(kns.Context)
	ccm, err := k.GetConfigMapWithKns(ctx, kns, ncm.Name, filter)
	switch {
	case ccm != nil && len(ccm.Name) > 0:
		switch IsVersionNew(ccm.Labels, ncm.Labels) {
		case true:
			derr := k.DeleteConfigMapWithKns(ctx, kns, ccm.Name, filter)
			if derr != nil {
				return ccm, derr
			}
		case false:
			return ccm, nil
		}
	case errors.IsNotFound(err):
		newCm, newCmErr := k.CreateConfigMapWithKns(ctx, kns, ncm, filter)
		return newCm, newCmErr
	}
	newCm, newSErr := k.CreateConfigMapWithKns(ctx, kns, ncm, filter)
	return newCm, newSErr
}

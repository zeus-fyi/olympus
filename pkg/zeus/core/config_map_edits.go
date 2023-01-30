package zeus_core

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) ConfigMapKeySwap(ctx context.Context, kns zeus_common_types.CloudCtxNs, name, key1, key2 string, filter *string_utils.FilterOpts) (*v1.ConfigMap, error) {
	k.SetContext(kns.Context)
	cm, err := k.GetConfigMapWithKns(ctx, kns, name, filter)
	if err != nil {
		return nil, err
	}
	v, ok := cm.Data[key1]
	v2, ok2 := cm.Data[key2]
	m := make(map[string]string)
	m = cm.Data
	if ok && ok2 {
		m[key1] = v2
		m[key2] = v
	} else {
		log.Ctx(ctx).Warn().Msg("key not found")
		return nil, err
	}
	cm.Data = m
	cmOut, err := k.kc.CoreV1().ConfigMaps(kns.Namespace).Update(ctx, cm, metav1.UpdateOptions{})
	return cmOut, err
}

func (k *K8Util) ConfigMapOverwriteOrCreateFromKey(ctx context.Context, kns zeus_common_types.CloudCtxNs, name, keyToCopy, keyToSetOrCreateFromCopy string, filter *string_utils.FilterOpts) (*v1.ConfigMap, error) {
	k.SetContext(kns.Context)
	cm, err := k.GetConfigMapWithKns(ctx, kns, name, filter)
	if err != nil {
		return nil, err
	}

	vSrc, ok := cm.Data[keyToCopy]
	m := make(map[string]string)
	m = cm.Data
	if ok {
		m[keyToSetOrCreateFromCopy] = vSrc
	} else {
		log.Ctx(ctx).Warn().Msg("key not found")
		return nil, err
	}
	cm.Data = m
	cmOut, err := k.kc.CoreV1().ConfigMaps(kns.Namespace).Update(ctx, cm, metav1.UpdateOptions{})
	return cmOut, err
}

func (k *K8Util) ConfigMapOverwriteOrCreateNewKeys(ctx context.Context, kns zeus_common_types.CloudCtxNs, name string, newKeyValMap map[string]string, filter *string_utils.FilterOpts) (*v1.ConfigMap, error) {
	k.SetContext(kns.Context)
	cm, err := k.GetConfigMapWithKns(ctx, kns, name, filter)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	if len(cm.Data) <= 0 {
		cm.Data = m
	}
	m = cm.Data
	for key, newVal := range newKeyValMap {
		currentVal, ok := cm.Data[key]
		if ok {
			log.Info().Interface("originalKeyValue", currentVal).Interface("newKeyVal", newVal)
		} else {
			log.Ctx(ctx).Info().Interface("newKey", key).Interface("newKeyVal", newVal).Msg("Creating new key value")
		}
		cm.Data[key] = newVal
	}
	cm.Data = m
	cmOut, err := k.kc.CoreV1().ConfigMaps(kns.Namespace).Update(ctx, cm, metav1.UpdateOptions{})
	return cmOut, err
}

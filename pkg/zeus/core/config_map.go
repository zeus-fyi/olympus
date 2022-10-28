package zeus_core

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetConfigMapWithKns(ctx context.Context, kns KubeCtxNs, name string, filter *string_utils.FilterOpts) (*v1.ConfigMap, error) {
	return k.kc.CoreV1().ConfigMaps(kns.Namespace).Get(context.Background(), name, metav1.GetOptions{})
}

func (k *K8Util) CreateConfigMapWithKns(ctx context.Context, kns KubeCtxNs, cm *v1.ConfigMap, filter *string_utils.FilterOpts) (*v1.ConfigMap, error) {
	return k.kc.CoreV1().ConfigMaps(kns.Namespace).Create(context.Background(), cm, metav1.CreateOptions{})
}

func (k *K8Util) DeleteConfigMapWithKns(ctx context.Context, kns KubeCtxNs, name string, filter *string_utils.FilterOpts) error {
	return k.kc.CoreV1().ConfigMaps(kns.Namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

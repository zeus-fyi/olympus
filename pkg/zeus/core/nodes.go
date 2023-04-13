package zeus_core

import (
	"context"

	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) GetNodes(ctx context.Context, kns zeus_common_types.CloudCtxNs) (*v1.NodeList, error) {
	k.SetContext(kns.Context)
	return k.kc.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
}

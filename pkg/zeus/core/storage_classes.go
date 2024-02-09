package zeus_core

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8Util) CreateStorageClass(ctx context.Context, kubeCtxNs zeus_common_types.CloudCtxNs, sc *v1.StorageClass) (*v1.StorageClass, error) {
	k.SetContext(kubeCtxNs.Context)

	// Create the storage class
	createdSc, err := k.kc.StorageV1().StorageClasses().Create(ctx, sc, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create storage class: %v", err)
	}
	return createdSc, nil
}

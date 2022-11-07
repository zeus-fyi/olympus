package destroy_deploy

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DeployConfigMap(ctx context.Context) error {
	return d.postDestroyDeployTarget("configmap")
}

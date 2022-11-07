package destroy_deploy

import (
	"context"
)

func (d *UndeployTopologyActivity) DeployConfigMap(ctx context.Context) error {
	return d.postDestroyDeployTarget("configmap")
}

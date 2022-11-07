package destroy_deploy_activities

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DeployConfigMap(ctx context.Context) error {
	return d.postDestroyDeployTarget("configmap")
}

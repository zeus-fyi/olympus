package destroy_deploy_activities

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DestroyDeployConfigMap(ctx context.Context) error {
	return d.postDestroyDeployTarget("configmap")
}

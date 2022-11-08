package destroy_deploy_activities

import (
	"context"
)

func (d *DestroyDeployTopologyActivities) DestroyDeployConfigMap(ctx context.Context) error {
	return d.postDestroyDeployTarget("configmap")
}

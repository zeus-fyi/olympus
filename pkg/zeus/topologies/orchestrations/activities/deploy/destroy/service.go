package destroy_deploy_activities

import (
	"context"
)

func (d *DestroyDeployTopologyActivities) DestroyDeployService(ctx context.Context) error {
	return d.postDestroyDeployTarget("service")
}

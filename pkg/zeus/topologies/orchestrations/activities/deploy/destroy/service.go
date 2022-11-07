package destroy_deploy_activities

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DestroyDeployService(ctx context.Context) error {
	return d.postDestroyDeployTarget("service")
}

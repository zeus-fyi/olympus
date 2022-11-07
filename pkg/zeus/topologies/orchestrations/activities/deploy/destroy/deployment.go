package destroy_deploy_activities

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DestroyDeployDeployment(ctx context.Context) error {
	return d.postDestroyDeployTarget("deployment")
}

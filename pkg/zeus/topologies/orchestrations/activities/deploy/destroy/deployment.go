package destroy_deploy_activities

import (
	"context"
)

func (d *DestroyDeployTopologyActivities) DestroyDeployDeployment(ctx context.Context) error {
	return d.postDestroyDeployTarget("deployment")
}

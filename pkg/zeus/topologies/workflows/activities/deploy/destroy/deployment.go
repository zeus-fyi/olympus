package destroy_deploy

import (
	"context"
)

func (d *UndeployTopologyActivity) DeployDeployment(ctx context.Context) error {
	return d.postDestroyDeployTarget("deployment")
}

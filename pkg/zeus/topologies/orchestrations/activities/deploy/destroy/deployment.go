package destroy_deploy

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DeployDeployment(ctx context.Context) error {
	return d.postDestroyDeployTarget("deployment")
}

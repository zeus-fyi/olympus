package deploy_topology

import (
	"context"
)

func (d *DeployTopologyActivity) DeployDeployment(ctx context.Context) error {
	return d.postDeployTarget("deployment")
}

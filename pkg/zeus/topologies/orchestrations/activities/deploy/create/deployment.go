package deploy_topology_activities

import (
	"context"
)

func (d *DeployTopologyActivity) DeployDeployment(ctx context.Context) error {
	return d.postDeployTarget("deployment")
}

package deploy_topology_activities

import (
	"context"
)

func (d *DeployTopologyActivity) DeployService(ctx context.Context) error {
	return d.postDeployTarget("service")
}

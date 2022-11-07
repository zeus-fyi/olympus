package deploy_topology

import (
	"context"
)

func (d *DeployTopologyActivity) DeployService(ctx context.Context) error {
	return d.postDeployTarget("service")
}

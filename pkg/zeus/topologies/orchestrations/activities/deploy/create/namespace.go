package deploy_topology_activities

import (
	"context"
)

func (d *DeployTopologyActivity) CreateNamespace(ctx context.Context) error {
	return d.postDeployTarget("namespace")
}

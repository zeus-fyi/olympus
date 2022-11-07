package deploy_topology

import (
	"context"
)

func (d *DeployTopologyActivity) CreateNamespace(ctx context.Context) error {
	return d.postDeployTarget("namespace")
}

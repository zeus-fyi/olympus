package deploy_topology

import (
	"context"
)

func (d *DeployTopologyActivity) DeployIngress(ctx context.Context) error {
	return d.postDeployTarget("ingress")
}

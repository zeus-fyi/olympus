package deploy_topology_activities

import (
	"context"
)

func (d *DeployTopologyActivity) DeployIngress(ctx context.Context) error {
	return d.postDeployTarget("ingress")
}

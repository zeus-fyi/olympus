package deploy_topology_activities

import (
	"context"
)

func (d *DeployTopologyActivities) DeployIngress(ctx context.Context) error {
	return d.postDeployTarget("ingress")
}

package deploy_topology_activities

import (
	"context"
)

func (d *DeployTopologyActivities) DeployConfigMap(ctx context.Context) error {
	return d.postDeployTarget("configmap")
}

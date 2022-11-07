package deploy_topology_activities

import (
	"context"
)

func (d *DeployTopologyActivity) DeployConfigMap(ctx context.Context) error {
	return d.postDeployTarget("configmap")
}

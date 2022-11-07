package deploy_topology_activities

import "context"

func (d *DeployTopologyActivity) DeployStatefulSet(ctx context.Context) error {
	return d.postDeployTarget("statefulset")
}

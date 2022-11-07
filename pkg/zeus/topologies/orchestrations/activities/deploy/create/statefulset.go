package deploy_topology_activities

import "context"

func (d *DeployTopologyActivities) DeployStatefulSet(ctx context.Context) error {
	return d.postDeployTarget("statefulset")
}

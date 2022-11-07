package deploy_topology

import "context"

func (d *DeployTopologyActivity) DeployStatefulSet(ctx context.Context) error {
	return d.postDeployTarget("statefulset")
}

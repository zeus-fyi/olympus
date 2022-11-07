package destroy_deploy_activities

import "context"

func (d *DestroyDeployTopologyActivity) DestroyDeployStatefulSet(ctx context.Context) error {
	return d.postDestroyDeployTarget("statefulset")
}

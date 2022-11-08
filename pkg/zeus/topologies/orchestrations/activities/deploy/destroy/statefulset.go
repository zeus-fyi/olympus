package destroy_deploy_activities

import "context"

func (d *DestroyDeployTopologyActivities) DestroyDeployStatefulSet(ctx context.Context) error {
	return d.postDestroyDeployTarget("statefulset")
}

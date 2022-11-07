package destroy_deploy

import "context"

func (d *DestroyDeployTopologyActivity) DestroyDeployStatefulSet(ctx context.Context) error {
	return d.postDestroyDeployTarget("statefulset")
}

package destroy_deploy_activities

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DestroyNamespace(ctx context.Context) error {
	return d.postDestroyDeployTarget("namespace")
}

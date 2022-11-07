package destroy_deploy

import (
	"context"
)

func (d *UndeployTopologyActivity) DeployIngress(ctx context.Context) error {
	return d.postDestroyDeployTarget("ingress")
}

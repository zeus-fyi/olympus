package destroy_deploy

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DeployIngress(ctx context.Context) error {
	return d.postDestroyDeployTarget("ingress")
}

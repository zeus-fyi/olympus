package destroy_deploy

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DestroyDeployIngress(ctx context.Context) error {
	return d.postDestroyDeployTarget("ingress")
}

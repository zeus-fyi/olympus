package destroy_deploy_activities

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DestroyDeployIngress(ctx context.Context) error {
	return d.postDestroyDeployTarget("ingress")
}

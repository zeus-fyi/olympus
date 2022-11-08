package destroy_deploy_activities

import (
	"context"
)

func (d *DestroyDeployTopologyActivities) DestroyDeployIngress(ctx context.Context) error {
	return d.postDestroyDeployTarget("ingress")
}

package destroy_deploy_activities

import (
	"context"

	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/base_request"
)

func (d *DestroyDeployTopologyActivities) DestroyDeployIngress(ctx context.Context, params base_request.InternalDeploymentActionRequest) error {
	return d.postDestroyDeployTarget("ingress", params)
}

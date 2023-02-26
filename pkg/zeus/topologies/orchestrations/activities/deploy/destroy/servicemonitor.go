package destroy_deploy_activities

import (
	"context"

	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
)

func (d *DestroyDeployTopologyActivities) DestroyDeployServiceMonitor(ctx context.Context, params base_request.InternalDeploymentActionRequest) error {
	return d.postDestroyDeployTarget("servicemonitor", params)
}

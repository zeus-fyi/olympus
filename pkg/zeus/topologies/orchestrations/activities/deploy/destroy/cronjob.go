package destroy_deploy_activities

import (
	"context"

	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
)

func (d *DestroyDeployTopologyActivities) DestroyCronJob(ctx context.Context, params base_request.InternalDeploymentActionRequest) error {
	return d.postDestroyDeployTarget("cronjob", params)
}

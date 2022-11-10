package deploy_topology_activities

import (
	"context"

	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/base_request"
)

func (d *DeployTopologyActivities) CreateNamespace(ctx context.Context, params base_request.InternalDeploymentActionRequest) error {
	return d.postDeployTarget("namespace", params)
}

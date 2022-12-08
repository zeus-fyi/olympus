package deploy_topology_activities

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

func (d *DeployTopologyActivities) DeployClusterTopology(ctx context.Context, params zeus_req_types.TopologyDeployRequest) error {
	return d.postDeployClusterTopology(params)
}

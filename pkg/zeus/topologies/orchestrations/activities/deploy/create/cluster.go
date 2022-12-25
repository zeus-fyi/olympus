package deploy_topology_activities

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
)

func (d *DeployTopologyActivities) DeployClusterTopology(ctx context.Context, params zeus_req_types.TopologyDeployRequest, ou org_users.OrgUser) error {
	return d.postDeployClusterTopology(params, ou)
}

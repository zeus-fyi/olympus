package deploy

import (
	topology_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities"
)

type DeploymentActionRequest struct {
	topology_activities.TopologyActivityRequest
}

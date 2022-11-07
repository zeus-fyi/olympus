package base_request

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	topology_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities"
)

type InternalDeploymentActionRequest struct {
	topology_activities.TopologyActivity
	chart_workload.NativeK8s
}

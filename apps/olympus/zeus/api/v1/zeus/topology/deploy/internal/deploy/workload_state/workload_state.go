package workload_state

import (
	create_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/state"
)

type InternalWorkloadStatusUpdate struct {
	create_topology_deployment_status.DeploymentStatus
}

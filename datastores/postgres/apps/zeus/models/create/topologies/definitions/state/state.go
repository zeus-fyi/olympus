package create_topology_deployment_status

import "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"

type DeploymentStatus struct {
	topology_deployment_status.Status
}

func NewCreateState() DeploymentStatus {
	s := topology_deployment_status.NewTopologyStatus()
	cs := DeploymentStatus{s}
	return cs
}

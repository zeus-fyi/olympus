package topology_deployment_status

import autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"

// Status specific deployed topology to user (statuses can be pending, terminated, etc)
type Status struct {
	autogen_bases.TopologiesDeployed
}

const DeployPending = "DeployPending"
const DeployInProgress = "DeployInProgress"
const DeployComplete = "DeployComplete"

const DestroyDeployPending = "DestroyDeployPending"
const DestroyDeployInProgress = "DestroyDeployInProgress"
const DestroyDeployComplete = "DestroyDeployComplete"

const InProgress = "InProgress"
const Complete = "Complete"

func NewPopulatedTopologyStatus(topID int, status string) Status {
	s := Status{autogen_bases.TopologiesDeployed{
		TopologyStatus: status,
		TopologyID:     topID,
	}}
	return s
}

func NewTopologyStatus() Status {
	s := Status{autogen_bases.TopologiesDeployed{
		TopologyStatus: "Pending",
		TopologyID:     0,
	}}
	return s
}

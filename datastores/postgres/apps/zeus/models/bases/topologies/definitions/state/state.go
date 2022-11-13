package topology_deployment_status

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
)

// Status specific deployed topology to user (statuses can be pending, terminated, etc)
type Status struct {
	kns.TopologyKubeCtxNs
	DeployStatus
}

type DeployStatus struct {
	autogen_bases.TopologiesDeployed
}

func (s *Status) SetTopologyID(id int) {
	s.DeployStatus.TopologyID = id
	s.TopologyKubeCtxNs.TopologyID = id
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
	s := Status{
		TopologyKubeCtxNs: kns.NewKns(),
		DeployStatus:      DeployStatus{autogen_bases.TopologiesDeployed{TopologyStatus: status, TopologyID: topID}},
	}
	s.SetTopologyID(topID)
	return s
}

func NewTopologyStatus() Status {
	s := Status{
		kns.NewKns(),
		DeployStatus{autogen_bases.TopologiesDeployed{
			TopologyStatus: "Pending",
			TopologyID:     0,
		}}}
	return s
}

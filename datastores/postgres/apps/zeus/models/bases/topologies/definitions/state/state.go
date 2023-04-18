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

type ClusterStatus struct {
	ClusterName string `json:"clusterName"`
	Status      string `json:"status"`
}

type DeployStatus struct {
	autogen_bases.TopologiesDeployed
}

func (s *Status) SetTopologyID(id int) {
	s.DeployStatus.TopologyID = id
	s.TopologyKubeCtxNs.TopologyID = id
}

const ResourceDestroyPending = "ResourceDestroyPending"

const DeployPending = "DeployPending"
const DeployInProgress = "DeployInProgress"
const DeployComplete = "DeployComplete"

const DestroyDeployPending = "DestroyDeployPending"
const DestroyDeployInProgress = "DestroyDeployInProgress"
const DestroyDeployComplete = "DestroyDeployComplete"

const CleanDeployPending = "CleanDeployPending"
const CleanDeployInProgress = "CleanDeployInProgress"
const CleanDeployComplete = "CleanDeployComplete"

const InProgress = "InProgress"
const Complete = "Complete"

func NewPopulatedTopologyStatus(kns kns.TopologyKubeCtxNs, status string) Status {
	s := Status{
		TopologyKubeCtxNs: kns,
		DeployStatus: DeployStatus{
			autogen_bases.TopologiesDeployed{
				TopologyStatus: status, TopologyID: kns.TopologyID},
		},
	}
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

func NewClusterTopologyStatus(name string) ClusterStatus {
	c := ClusterStatus{
		ClusterName: name,
		Status:      "Pending",
	}
	return c
}

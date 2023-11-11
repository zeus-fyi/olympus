package topology_deployment_status

import (
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

// Status specific deployed topology to user (statuses can be pending, terminated, etc)
type Status struct {
	zeus_req_types.TopologyDeployRequest `json:"topologyDeployRequest"`
	DeployStatus                         `json:"deployStatus"`
}

type ClusterStatus struct {
	ClusterName string `json:"clusterName"`
	Status      string `json:"status"`
}

type DeployStatus struct {
	autogen_bases.TopologiesDeployed `json:"deployStatus"`
}

func (s *Status) SetTopologyID(id int) {
	s.DeployStatus.TopologyID = id
	s.TopologyDeployRequest.TopologyID = id
}

const (
	ResourceDestroyPending           = "ResourceDestroyPending"
	DeployPending                    = "DeployPending"
	FleetRolloutRestartDeployPending = "FleetRolloutRestartDeployPending"
	DeployInProgress                 = "DeployInProgress"
	DeployComplete                   = "DeployComplete"

	DestroyDeployPending    = "DestroyDeployPending"
	DestroyDeployInProgress = "DestroyDeployInProgress"
	DestroyDeployComplete   = "DestroyDeployComplete"

	CleanDeployPending    = "CleanDeployPending"
	CleanDeployInProgress = "CleanDeployInProgress"
	CleanDeployComplete   = "CleanDeployComplete"

	InProgress = "InProgress"
	Complete   = "Complete"
)

func NewPopulatedTopologyStatus(kns zeus_req_types.TopologyDeployRequest, status string) Status {
	s := Status{
		TopologyDeployRequest: kns,
		DeployStatus: DeployStatus{
			autogen_bases.TopologiesDeployed{
				TopologyStatus: status, TopologyID: kns.TopologyID},
		},
	}
	return s
}

func NewTopologyStatus() Status {
	s := Status{
		TopologyDeployRequest: zeus_req_types.TopologyDeployRequest{
			TopologyID:                      0,
			ClusterClassName:                "",
			CloudCtxNs:                      zeus_common_types.CloudCtxNs{},
			SecretRef:                       "",
			RequestChoreographySecretDeploy: false,
		},
		DeployStatus: DeployStatus{
			autogen_bases.TopologiesDeployed{
				TopologyID:     0,
				TopologyStatus: "Pending",
			},
		},
	}
	return s
}

func NewClusterTopologyStatus(name string) ClusterStatus {
	c := ClusterStatus{
		ClusterName: name,
		Status:      "Pending",
	}
	return c
}

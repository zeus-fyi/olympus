package base_deploy_params

import (
	"net/url"
	"path"

	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

type TopologyWorkflowRequest struct {
	TopologyDeployRequest zeus_req_types.TopologyDeployRequest `json:"topologyDeployRequest"`
	OrgUser               org_users.OrgUser                    `json:"orgUser"`
}

type FleetUpgradeWorkflowRequest struct {
	OrgUser     org_users.OrgUser `json:"orgUser"`
	Host        string            `json:"host,omitempty"`
	ClusterName string            `json:"clusterName"`
	AppTaint    bool              `json:"appTaint"`
}

type FleetRolloutRestartWorkflowRequest struct {
	OrgUser     org_users.OrgUser `json:"orgUser"`
	ClusterName string            `json:"clusterName"`
}

type ClusterTopologyWorkflowRequest struct {
	ClusterClassName string                       `json:"clusterClassName"`
	TopologyIDs      []int                        `json:"topologyIDs"`
	CloudCtxNs       zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
	OrgUser          org_users.OrgUser            `json:"orgUser"`
	Host             string                       `json:"host,omitempty"`
	AppTaint         bool                         `json:"appTaint"`

	RequestChoreographySecret bool `json:"requestChoreographySecret"`
}

type DestroyResourcesRequest struct {
	Ou             org_users.OrgUser `json:"orgUser"`
	OrgResourceIDs []int             `json:"orgResourceIDs"`
}

type DestroyClusterSetupRequest struct {
	ClusterSetupRequest `json:"clusterSetupRequest"`
}

type ClusterSetupRequest struct {
	FreeTrial     bool                            `json:"freeTrial"`
	Ou            org_users.OrgUser               `json:"orgUser"`
	CloudCtxNs    zeus_common_types.CloudCtxNs    `json:"cloudCtxNs"`
	Nodes         hestia_autogen_bases.Nodes      `json:"nodes"`
	NodesQuantity float64                         `json:"nodesQuantity"`
	Disks         hestia_autogen_bases.DisksSlice `json:"disks"`
	Cluster       zeus_templates.Cluster          `json:"cluster"`
	AppTaint      bool                            `json:"appTaint"`
}

func (t *TopologyWorkflowRequest) GetURL(prefix, target string) url.URL {
	u := url.URL{
		Host: "https://api.zeus.fyi",
		Path: path.Join(prefix, target),
	}
	return u
}

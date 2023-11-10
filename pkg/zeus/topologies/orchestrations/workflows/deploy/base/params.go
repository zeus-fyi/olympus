package base_deploy_params

import (
	"net/url"
	"path"

	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type TopologyWorkflowRequest struct {
	Kns                       kns.TopologyKubeCtxNs
	OrgUser                   org_users.OrgUser
	Host                      string `json:"host,omitempty"`
	RequestChoreographySecret bool   `json:"requestChoreographySecret"`
	ClusterClassName          string `json:"clusterClassName"`
	SecretRef                 string `json:"secretRef"`

	chart_workload.TopologyBaseInfraWorkload
}

type FleetUpgradeWorkflowRequest struct {
	OrgUser     org_users.OrgUser `json:"orgUser"`
	Host        string            `json:"host,omitempty"`
	ClusterName string            `json:"clusterName"`
	AppTaint    bool              `json:"appTaint"`
}

type ClusterTopologyWorkflowRequest struct {
	ClusterClassName string                       `json:"clusterClassName"`
	TopologyIDs      []int                        `json:"topologyIDs"`
	CloudCtxNS       zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
	OrgUser          org_users.OrgUser            `json:"orgUser"`
	Host             string                       `json:"host,omitempty"`
	AppTaint         bool                         `json:"appTaint"`

	RequestChoreographySecret bool `json:"requestChoreographySecret"`
}

type DestroyResourcesRequest struct {
	Ou             org_users.OrgUser
	OrgResourceIDs []int
}

type DestroyClusterSetupRequest struct {
	ClusterSetupRequest
}

type ClusterSetupRequest struct {
	FreeTrial                    bool              `json:"freeTrial"`
	Ou                           org_users.OrgUser `json:"orgUser"`
	zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
	Nodes                        hestia_autogen_bases.Nodes      `json:"nodes"`
	NodesQuantity                float64                         `json:"nodesQuantity"`
	Disks                        hestia_autogen_bases.DisksSlice `json:"disks"`
	Cluster                      zeus_templates.Cluster          `json:"cluster"`
	AppTaint                     bool                            `json:"appTaint"`
}

func (t *TopologyWorkflowRequest) GetURL(prefix, target string) url.URL {
	if len(t.Host) <= 0 {
		t.Host = "https://api.zeus.fyi"
	}
	u := url.URL{
		Host: t.Host,
		Path: path.Join(prefix, target),
	}
	return u
}

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
	Host                      string
	RequestChoreographySecret bool
	ClusterClassName          string `json:"clusterClassName"`
	SecretRef                 string `json:"secretRef"`

	chart_workload.TopologyBaseInfraWorkload
}

type FleetUpgradeWorkflowRequest struct {
	OrgUser     org_users.OrgUser
	Host        string
	ClusterName string
	AppTaint    bool
}

type ClusterTopologyWorkflowRequest struct {
	ClusterClassName string `json:"clusterClassName"`
	TopologyIDs      []int
	CloudCtxNS       zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
	OrgUser          org_users.OrgUser
	Host             string
	AppTaint         bool

	RequestChoreographySecret bool
}

type DestroyResourcesRequest struct {
	Ou             org_users.OrgUser
	OrgResourceIDs []int
}

type DestroyClusterSetupRequest struct {
	ClusterSetupRequest
}

type ClusterSetupRequest struct {
	FreeTrial bool
	Ou        org_users.OrgUser
	zeus_common_types.CloudCtxNs
	Nodes         hestia_autogen_bases.Nodes
	NodesQuantity float64
	Disks         hestia_autogen_bases.DisksSlice
	Cluster       zeus_templates.Cluster
	AppTaint      bool
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

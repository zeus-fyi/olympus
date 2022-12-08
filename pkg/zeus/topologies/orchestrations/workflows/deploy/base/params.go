package base_deploy_params

import (
	"net/url"
	"path"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

type TopologyWorkflowRequest struct {
	Kns     kns.TopologyKubeCtxNs
	OrgUser org_users.OrgUser
	Host    string

	chart_workload.TopologyBaseInfraWorkload
}

type ClusterTopologyWorkflowRequest struct {
	ClusterName string
	TopologyIDs []int
	CloudCtxNS  zeus_common_types.CloudCtxNs
	OrgUser     org_users.OrgUser
	Host        string
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

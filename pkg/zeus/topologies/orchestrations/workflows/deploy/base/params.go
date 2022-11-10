package base_deploy_params

import (
	"net/url"
	"path"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
)

type TopologyWorkflowRequest struct {
	Kns     kns.TopologyKubeCtxNs
	OrgUser org_users.OrgUser
	Bearer  string
	Host    string

	chart_workload.NativeK8s
}

func (t *TopologyWorkflowRequest) GetURL(target string) url.URL {
	if len(t.Host) <= 0 {
		t.Host = "https://api.zeus.fyi"
	}
	u := url.URL{
		Host: t.Host,
		Path: path.Join(t.Host, target),
	}
	return u
}

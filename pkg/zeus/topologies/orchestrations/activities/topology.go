package activities

import (
	"net/url"
	"path"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type TopologyActivity struct {
	temporal_base.Activity
	chart_workload.NativeK8s
	OrgUser org_users.OrgUser
	Kns     zeus_core.KubeCtxNs
	Host    string
}

func (t *TopologyActivity) GetURL(prefix, target string) url.URL {
	if len(t.Host) <= 0 {
		t.Host = "https://api.zeus.fyi"
	}
	u := url.URL{
		Host: t.Host,
		Path: path.Join(prefix, target),
	}
	return u
}

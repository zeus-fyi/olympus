package base_request

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
)

type InternalDeploymentActionRequest struct {
	Kns     kns.TopologyKubeCtxNs
	OrgUser org_users.OrgUser
	chart_workload.TopologyBaseInfraWorkload
}

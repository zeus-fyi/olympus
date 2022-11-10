package helpers

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

func PackageCommonTopologyRequest(topCtxNs kns.TopologyKubeCtxNs, ou org_users.OrgUser, nk chart_workload.NativeK8s) base_deploy_params.TopologyWorkflowRequest {
	tar := base_deploy_params.TopologyWorkflowRequest{
		Kns:       topCtxNs,
		OrgUser:   ou,
		Host:      "",
		NativeK8s: nk,
	}
	return tar
}

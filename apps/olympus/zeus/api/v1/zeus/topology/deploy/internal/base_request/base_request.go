package base_request

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	deploy_topology "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy"
)

type InternalDeploymentActionRequest struct {
	deploy_topology.DeployParams
	chart_workload.NativeK8s
}

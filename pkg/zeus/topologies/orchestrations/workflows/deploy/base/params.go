package base_deploy_params

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type DeployTopologyParams struct {
	Kns        zeus_core.KubeCtxNs
	TopologyID int
	UserID     int
	OrgID      int
	chart_workload.NativeK8s
}

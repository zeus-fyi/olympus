package base_request

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type InternalDeploymentActionRequest struct {
	TopologyActivity
	chart_workload.NativeK8s
}

type TopologyActivity struct {
	temporal_base.Activity
	chart_workload.NativeK8s

	Bearer string
	Kns    zeus_core.KubeCtxNs
	Host   string
}

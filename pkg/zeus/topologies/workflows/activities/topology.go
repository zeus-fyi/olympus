package activities

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type TopologyActivity struct {
	zeus_core.K8Util
	chart_workload.NativeK8s
	Kns zeus_core.KubeCtxNs
}

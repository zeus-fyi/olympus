package apollo_metrics_workload_info

import "github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"

type WorkloadInfo struct {
	Namespace string
	Subsystem string
	zeus_common_types.CloudCtxNs
}

func NewWorkloadInfo(subsystem string, ctxNs zeus_common_types.CloudCtxNs) WorkloadInfo {
	return WorkloadInfo{
		Namespace:  ctxNs.Namespace,
		Subsystem:  subsystem,
		CloudCtxNs: ctxNs,
	}
}

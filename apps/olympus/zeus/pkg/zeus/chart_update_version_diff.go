package zeus

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

func DiffChartUpdate(replacement, current chart_workload.TopologyBaseInfraWorkload) chart_workload.TopologyBaseInfraWorkload {
	ts := chronos.Chronos{}
	workloadDiff := chart_workload.TopologyBaseInfraWorkload{}
	if replacement.ConfigMap != nil && current.ConfigMap != nil && replacement.ConfigMap.Name == current.ConfigMap.Name {
		workloadDiff.ConfigMap = replacement.ConfigMap
		workloadDiff.ConfigMap.Labels["version"] = fmt.Sprintf("%d", ts.UnixTimeStampNow())
	}
	if replacement.StatefulSet != nil && current.StatefulSet != nil && replacement.StatefulSet.Name == current.StatefulSet.Name {
		workloadDiff.StatefulSet = replacement.StatefulSet
		workloadDiff.StatefulSet.Labels["version"] = fmt.Sprintf("%d", ts.UnixTimeStampNow())
	}
	if replacement.Deployment != nil && current.Deployment != nil && replacement.Deployment.Name == current.Deployment.Name {
		workloadDiff.Deployment = replacement.Deployment
		workloadDiff.Deployment.Labels["version"] = fmt.Sprintf("%d", ts.UnixTimeStampNow())
	}
	if replacement.Service != nil && current.Service != nil && replacement.Service.Name == current.Service.Name {
		workloadDiff.Service = replacement.Service
		workloadDiff.Service.Labels["version"] = fmt.Sprintf("%d", ts.UnixTimeStampNow())
	}
	if replacement.Ingress != nil && current.Ingress != nil && replacement.Ingress.Name == current.Ingress.Name {
		workloadDiff.Ingress = replacement.Ingress
		workloadDiff.Ingress.Labels["version"] = fmt.Sprintf("%d", ts.UnixTimeStampNow())
	}
	return workloadDiff

}

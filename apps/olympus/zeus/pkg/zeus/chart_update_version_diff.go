package zeus

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

func DiffChartUpdate(replacement, current chart_workload.TopologyBaseInfraWorkload) chart_workload.TopologyBaseInfraWorkload {
	workloadDiff := chart_workload.TopologyBaseInfraWorkload{}
	if replacement.ConfigMap != nil && current.ConfigMap != nil && replacement.ConfigMap.Name == current.ConfigMap.Name {
		workloadDiff.ConfigMap = replacement.ConfigMap
		workloadDiff.ConfigMap.Labels = UpdateVersionIDLabel(workloadDiff.ConfigMap.Labels)
	}
	if replacement.StatefulSet != nil && current.StatefulSet != nil && replacement.StatefulSet.Name == current.StatefulSet.Name {
		workloadDiff.StatefulSet = replacement.StatefulSet
		workloadDiff.StatefulSet.Labels = UpdateVersionIDLabel(workloadDiff.StatefulSet.Labels)
	}
	if replacement.Deployment != nil && current.Deployment != nil && replacement.Deployment.Name == current.Deployment.Name {
		workloadDiff.Deployment = replacement.Deployment
		workloadDiff.Deployment.Labels = UpdateVersionIDLabel(workloadDiff.Deployment.Labels)
	}
	if replacement.Service != nil && current.Service != nil && replacement.Service.Name == current.Service.Name {
		workloadDiff.Service = replacement.Service
		workloadDiff.Service.Labels = UpdateVersionIDLabel(workloadDiff.Service.Labels)
	}
	if replacement.Ingress != nil && current.Ingress != nil && replacement.Ingress.Name == current.Ingress.Name {
		workloadDiff.Ingress = replacement.Ingress
		workloadDiff.Ingress.Labels = UpdateVersionIDLabel(workloadDiff.Ingress.Labels)
	}
	return workloadDiff
}

func UpdateVersionIDLabel(labels map[string]string) map[string]string {
	if len(labels) <= 0 {
		labels = make(map[string]string)
	}
	var ts chronos.Chronos
	labels["version"] = fmt.Sprintf("version-%d", ts.UnixTimeStampNow())
	return labels
}

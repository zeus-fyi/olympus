package pods

import (
	"time"

	v1 "k8s.io/api/core/v1"
)

func parseResp(pods *v1.PodList) PodsSummary {
	if pods == nil {
		return PodsSummary{}
	}
	ps := PodsSummary{Pods: make(map[string]PodSummary, len(pods.Items))}
	for _, pod := range pods.Items {

		var ts time.Time
		if pod.Status.StartTime != nil {
			ts = pod.Status.StartTime.Time
		}
		podSum := PodSummary{
			PodName:       pod.GetName(),
			Phase:         string(pod.Status.Phase),
			Message:       pod.Status.Message,
			Reason:        pod.Status.Reason,
			StartTime:     ts,
			PodConditions: pod.Status.Conditions,
		}
		initContStatuses := map[string]v1.ContainerStatus{}
		for _, ic := range pod.Status.InitContainerStatuses {
			initContStatuses[ic.Name] = ic
		}

		contStatuses := map[string]v1.ContainerStatus{}
		for _, cont := range pod.Status.ContainerStatuses {
			contStatuses[cont.Name] = cont
		}
		cs := make([]string, len(pod.Spec.Containers))
		for i, cont := range pod.Spec.Containers {
			cs[i] = cont.Name
		}
		podSum.Containers = cs
		ps.Pods[pod.GetName()] = podSum
	}
	return ps
}

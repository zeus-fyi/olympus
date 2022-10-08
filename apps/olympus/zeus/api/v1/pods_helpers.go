package v1

import v1 "k8s.io/api/core/v1"

func parseResp(pods *v1.PodList) PodsSummary {
	if pods == nil {
		return PodsSummary{}
	}
	ps := PodsSummary{Pods: make(map[string]PodSummary, len(pods.Items))}
	for _, pod := range pods.Items {
		podSum := PodSummary{
			PodName:       pod.GetName(),
			Phase:         string(pod.Status.Phase),
			Message:       pod.Status.Message,
			Reason:        pod.Status.Reason,
			StartTime:     pod.Status.StartTime.Time,
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
		ps.Pods[pod.GetName()] = podSum
	}
	return ps
}

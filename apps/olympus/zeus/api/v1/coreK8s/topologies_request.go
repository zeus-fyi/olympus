package coreK8s

type TopologyActionRequest struct {
	K8sRequest
	Action        string
	ContainerName string
}

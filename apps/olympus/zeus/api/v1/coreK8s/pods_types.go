package coreK8s

import (
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	v12 "github.com/zeus-fyi/olympus/zeus/api/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodActionRequest struct {
	v12.K8sRequest
	Action        string
	PodName       string
	ContainerName string

	FilterOpts *string_utils.FilterOpts
	ClientReq  *ClientRequest
	LogOpts    *v1.PodLogOptions
	DeleteOpts *metav1.DeleteOptions
}

type ClientRequest struct {
	MethodHTTP      string
	Endpoint        string
	Ports           []string
	Payload         *string
	PayloadBytes    *[]byte
	EndpointHeaders map[string]string
}

type ClientResp struct {
	ReplyBodies map[string]string
}

type PodsSummary struct {
	Pods map[string]PodSummary `json:"pods"`
}

type PodSummary struct {
	PodName               string                        `json:"podName"`
	Phase                 string                        `json:"podPhase"`
	Message               string                        `json:"message"`
	Reason                string                        `json:"reason"`
	StartTime             time.Time                     `json:"startTime"`
	PodConditions         []v1.PodCondition             `json:"podConditions"`
	InitContainerStatuses map[string]v1.ContainerStatus `json:"initContainerConditions"`
	ContainerStatuses     map[string]v1.ContainerStatus `json:"containerStatuses"`
}

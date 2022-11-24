package zeus_pods_reqs

import (
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodActionRequest struct {
	zeus_req_types.TopologyDeployRequest
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

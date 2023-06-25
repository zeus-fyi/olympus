package hestia_ovhcloud

import (
	"fmt"
	"time"
)

type OvhNodePoolCreationRequest struct {
	ServiceName                 string `json:"serviceName"`
	KubeId                      string `json:"kubeId"`
	ProjectKubeNodePoolCreation `json:"ProjectKubeNodePoolCreation"`
}

func (op *OvhNodePoolCreationRequest) GetEndpoint() string {
	return fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", op.ServiceName, op.KubeId)
}

type ProjectKubeNodePoolCreation struct {
	AntiAffinity  *bool         `json:"antiAffinity,omitempty"`
	Autoscale     *bool         `json:"autoscale,omitempty"`
	Autoscaling   *Autoscaling  `json:"autoscaling,omitempty"`
	DesiredNodes  int           `json:"desiredNodes,omitempty"`
	FlavorName    string        `json:"flavorName"`
	MaxNodes      int           `json:"maxNodes,omitempty"`
	MinNodes      int           `json:"minNodes,omitempty"`
	MonthlyBilled *bool         `json:"monthlyBilled,omitempty"`
	Name          string        `json:"name,omitempty"`
	Template      *NodeTemplate `json:"template,omitempty"`
}

type Autoscaling struct {
	ScaleDownUnneededTimeSeconds  int     `json:"scaleDownUnneededTimeSeconds,omitempty"`
	ScaleDownUnreadyTimeSeconds   int     `json:"scaleDownUnreadyTimeSeconds,omitempty"`
	ScaleDownUtilizationThreshold float64 `json:"scaleDownUtilizationThreshold,omitempty"`
}

type NodeTemplate struct {
	Metadata *Metadata `json:"metadata,omitempty"`
	Spec     *Spec     `json:"spec,omitempty"`
}

type Metadata struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Finalizers  []string          `json:"finalizers,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type Spec struct {
	Taints        []KubernetesTaint `json:"taints,omitempty"`
	Unschedulable bool              `json:"unschedulable"`
}

type KubernetesTaint struct {
	Effect string `json:"effect"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type OvhNodePoolCreationResponse struct {
	AntiAffinity   bool          `json:"antiAffinity"`
	Autoscale      bool          `json:"autoscale"`
	Autoscaling    Autoscaling   `json:"autoscaling"`
	AvailableNodes int           `json:"availableNodes"`
	CreatedAt      time.Time     `json:"createdAt"`
	CurrentNodes   int           `json:"currentNodes"`
	DesiredNodes   int           `json:"desiredNodes"`
	Flavor         string        `json:"flavor"`
	Id             string        `json:"id"`
	MaxNodes       int           `json:"maxNodes"`
	MinNodes       int           `json:"minNodes"`
	MonthlyBilled  bool          `json:"monthlyBilled"`
	Name           string        `json:"name"`
	ProjectId      string        `json:"projectId"`
	SizeStatus     string        `json:"sizeStatus"`
	Status         string        `json:"status"`
	Template       *NodeTemplate `json:"template,omitempty"`
}

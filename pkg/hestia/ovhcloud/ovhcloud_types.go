package hestia_ovhcloud

import (
	"fmt"
	"time"
)

type OvhNodePoolCreationRequest struct {
	ServiceName                 string `json:"serviceName"`
	KubeId                      string `json:"kubeID"`
	ProjectKubeNodePoolCreation `json:"projectKubeNodePoolCreation"`
}

func (op *OvhNodePoolCreationRequest) GetEndpoint() string {
	return fmt.Sprintf("/cloud/project/%s}/kube/%s/nodepool", op.ServiceName, op.KubeId)
}

type ProjectKubeNodePoolCreation struct {
	AntiAffinity bool `json:"antiAffinity,omitempty"`
	Autoscale    bool `json:"autoscale,omitempty"`
	Autoscaling  struct {
		ScaleDownUnneededTimeSeconds  int `json:"scaleDownUnneededTimeSeconds"`
		ScaleDownUnreadyTimeSeconds   int `json:"scaleDownUnreadyTimeSeconds"`
		ScaleDownUtilizationThreshold int `json:"scaleDownUtilizationThreshold"`
	} `json:"autoscaling"`
	DesiredNodes  int    `json:"desiredNodes"`
	FlavorName    string `json:"flavorName"`
	MaxNodes      int    `json:"maxNodes"`
	MinNodes      int    `json:"minNodes"`
	MonthlyBilled bool   `json:"monthlyBilled,omitempty"`
	Name          string `json:"name"`
	Template      struct {
		Metadata struct {
			Annotations struct {
			} `json:"annotations"`
			Finalizers []string `json:"finalizers"`
			Labels     struct {
			} `json:"labels"`
		} `json:"metadata"`
		Spec struct {
			Taints []struct {
				Effect string `json:"effect"`
				Key    string `json:"key"`
				Value  string `json:"value"`
			} `json:"taints"`
			Unschedulable bool `json:"unschedulable"`
		} `json:"spec"`
	} `json:"template"`
}

type OvhNodePoolCreationResponse struct {
	AntiAffinity bool `json:"antiAffinity"`
	Autoscale    bool `json:"autoscale"`
	Autoscaling  struct {
		ScaleDownUnneededTimeSeconds  int `json:"scaleDownUnneededTimeSeconds"`
		ScaleDownUnreadyTimeSeconds   int `json:"scaleDownUnreadyTimeSeconds"`
		ScaleDownUtilizationThreshold int `json:"scaleDownUtilizationThreshold"`
	} `json:"autoscaling"`
	AvailableNodes int       `json:"availableNodes"`
	CreatedAt      time.Time `json:"createdAt"`
	CurrentNodes   int       `json:"currentNodes"`
	DesiredNodes   int       `json:"desiredNodes"`
	Flavor         string    `json:"flavor"`
	Id             string    `json:"id"`
	MaxNodes       int       `json:"maxNodes"`
	MinNodes       int       `json:"minNodes"`
	MonthlyBilled  bool      `json:"monthlyBilled"`
	Name           string    `json:"name"`
	ProjectId      string    `json:"projectId"`
	SizeStatus     string    `json:"sizeStatus"`
	Status         string    `json:"status"`
	Template       struct {
		Metadata struct {
			Annotations struct {
			} `json:"annotations"`
			Finalizers []string `json:"finalizers"`
			Labels     struct {
			} `json:"labels"`
		} `json:"metadata"`
		Spec struct {
			Taints []struct {
				Effect string `json:"effect"`
				Key    string `json:"key"`
				Value  string `json:"value"`
			} `json:"taints"`
			Unschedulable bool `json:"unschedulable"`
		} `json:"spec"`
	} `json:"template"`
	UpToDateNodes int       `json:"upToDateNodes"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

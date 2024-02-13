package zeus_templates

import (
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/zeus/zeus/workload_config_drivers/zk8s_templates"
)

type Cluster struct {
	ClusterName     string                 `json:"clusterName"`
	ComponentBases  ComponentBases         `json:"componentBases"`
	IngressSettings Ingress                `json:"ingressSettings"`
	IngressPaths    map[string]IngressPath `json:"ingressPaths"`
}

type IngressPaths map[string]IngressPath

type ComponentBases map[string]SkeletonBases

type SkeletonBases map[string]SkeletonBase

type SkeletonBase struct {
	TopologyID        string      `json:"topologyID,omitempty"`
	AddStatefulSet    bool        `json:"addStatefulSet"`
	AddDeployment     bool        `json:"addDeployment"`
	AddConfigMap      bool        `json:"addConfigMap"`
	AddService        bool        `json:"addService"`
	AddIngress        bool        `json:"addIngress"`
	AddServiceMonitor bool        `json:"addServiceMonitor"`
	ConfigMap         ConfigMap   `json:"configMap"`
	Deployment        Deployment  `json:"deployment"`
	StatefulSet       StatefulSet `json:"statefulSet"`
	Containers        Containers  `json:"containers"`

	ResourceSums zeus_core.ResourceSums `json:"resourceSums,omitempty"`
}

type ConfigMap map[string]string

type Ingress struct {
	AuthServerURL string `json:"authServerURL"`
	Host          string `json:"host"`
}

type IngressPath struct {
	Path     string `json:"path"`
	PathType string `json:"pathType"`
}

type Deployment struct {
	ReplicaCount int `json:"replicaCount"`
}

type StatefulSet struct {
	ReplicaCount int                          `json:"replicaCount"`
	PVCTemplates []zk8s_templates.PVCTemplate `json:"pvcTemplates"`
}

type Containers map[string]Container

type Container struct {
	IsInitContainer bool        `json:"isInitContainer"`
	ImagePullPolicy string      `json:"imagePullPolicy,omitempty"`
	DockerImage     DockerImage `json:"dockerImage"`
}

type DockerImage struct {
	ImageName            string               `json:"imageName"`
	Cmd                  string               `json:"cmd"`
	Args                 string               `json:"args"`
	ResourceRequirements ResourceRequirements `json:"resourceRequirements,omitempty"`
	Ports                []Port               `json:"ports,omitempty"`
	VolumeMounts         []VolumeMount        `json:"volumeMounts,omitempty"`
}

type ResourceRequirements struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type Port struct {
	Name               string `json:"name"`
	Number             string `json:"number"`
	Protocol           string `json:"protocol"`
	IngressEnabledPort bool   `json:"ingressEnabledPort"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

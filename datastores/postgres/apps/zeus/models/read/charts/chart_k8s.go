package read_charts

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/configuration"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulset"
	v1 "k8s.io/api/apps/v1"
	v1core "k8s.io/api/core/v1"
	v1networking "k8s.io/api/networking/v1"
)

type ChartWorkload struct {
	*deployments.Deployment
	*services.Service
	*ingresses.Ingress
	*configuration.ConfigMap
	*statefulset.StatefulSet
}

func NewK8sWorkload() ChartWorkload {
	k8s := ChartWorkload{
		Deployment: nil,
		Service:    nil,
		Ingress:    nil,
		ConfigMap:  nil,
	}
	return k8s
}

type NativeK8s struct {
	*v1core.Service       `json:"service"`
	*v1core.ConfigMap     `json:"configMap"`
	*v1.Deployment        `json:"deployment"`
	*v1.StatefulSet       `json:"statefulSet"`
	*v1networking.Ingress `json:"ingress"`
}

func (c *ChartWorkload) GetNativeK8s() NativeK8s {
	nk := NativeK8s{}

	if c.Deployment != nil {
		nk.Deployment = &c.K8sDeployment
	}
	if c.StatefulSet != nil {
		nk.StatefulSet = &c.K8sStatefulSet
	}
	if c.Service != nil {
		nk.Service = &c.K8sService
	}
	if c.ConfigMap != nil {
		nk.ConfigMap = &c.K8sConfigMap
	}
	if c.Ingress != nil {
		nk.Ingress = &c.K8sIngress
	}
	return nk
}

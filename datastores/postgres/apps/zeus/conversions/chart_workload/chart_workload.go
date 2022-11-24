package chart_workload

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/configuration"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulset"
)

type ChartWorkload struct {
	*deployments.Deployment
	*services.Service
	*ingresses.Ingress
	*configuration.ConfigMap
	*statefulset.StatefulSet
}

func NewChartWorkload() ChartWorkload {
	k8s := ChartWorkload{
		Deployment: nil,
		Service:    nil,
		Ingress:    nil,
		ConfigMap:  nil,
	}
	return k8s
}

func (c *ChartWorkload) GetTopologyBaseInfraWorkload() TopologyBaseInfraWorkload {
	nk := TopologyBaseInfraWorkload{}

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

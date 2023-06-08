package chart_workload

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/configuration"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/jobs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/servicemonitors"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulsets"
)

type ChartWorkload struct {
	*deployments.Deployment
	*services.Service
	*ingresses.Ingress
	*configuration.ConfigMap
	*statefulsets.StatefulSet
	*servicemonitors.ServiceMonitor
	*jobs.Job
	*jobs.CronJob
}

func NewChartWorkload() ChartWorkload {
	k8s := ChartWorkload{
		Deployment:     nil,
		Service:        nil,
		Ingress:        nil,
		ConfigMap:      nil,
		StatefulSet:    nil,
		ServiceMonitor: nil,
		Job:            nil,
		CronJob:        nil,
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
	if c.ServiceMonitor != nil {
		nk.ServiceMonitor = &c.K8sServiceMonitor
	}
	if c.Job != nil {
		nk.Job = &c.K8sJob
	}
	if c.CronJob != nil {
		nk.CronJob = &c.K8sCronJob
	}
	return nk
}

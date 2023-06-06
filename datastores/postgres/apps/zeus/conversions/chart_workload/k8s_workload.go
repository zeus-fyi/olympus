package chart_workload

import (
	v1sm "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/configuration"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/jobs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/ingresses"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/servicemonitors"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/statefulsets"
	v1 "k8s.io/api/apps/v1"
	v1Batch "k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	v1networking "k8s.io/api/networking/v1"
)

type TopologyBaseInfraWorkload struct {
	*v1Batch.CronJob      `json:"cronjob"`
	*v1Batch.Job          `json:"job"`
	*v1core.Service       `json:"service"`
	*v1core.ConfigMap     `json:"configMap"`
	*v1.Deployment        `json:"deployment"`
	*v1.StatefulSet       `json:"statefulSet"`
	*v1networking.Ingress `json:"ingress"`
	*v1sm.ServiceMonitor  `json:"serviceMonitor"`
}

func NewTopologyBaseInfraWorkload() TopologyBaseInfraWorkload {
	return TopologyBaseInfraWorkload{
		CronJob:        nil,
		Job:            nil,
		Service:        nil,
		ConfigMap:      nil,
		Deployment:     nil,
		StatefulSet:    nil,
		Ingress:        nil,
		ServiceMonitor: nil,
	}
}

func (nk *TopologyBaseInfraWorkload) CreateChartWorkloadFromTopologyBaseInfraWorkload() (ChartWorkload, error) {
	cw := NewChartWorkload()
	if nk.Deployment != nil {
		nd := deployments.NewDeployment()
		nd.K8sDeployment = *nk.Deployment
		err := nd.ConvertDeploymentConfigToDB()
		if err != nil {
			return cw, err
		}
		cw.Deployment = &nd
	}
	if nk.StatefulSet != nil {
		sts := statefulsets.NewStatefulSet()
		sts.K8sStatefulSet = *nk.StatefulSet
		err := sts.ConvertK8sStatefulSetToDB()
		if err != nil {
			return cw, err
		}
		cw.StatefulSet = &sts
	}
	if nk.Service != nil {
		nsvc := services.NewService()
		nsvc.K8sService = *nk.Service
		nsvc.ConvertK8sServiceToDB()
		cw.Service = &nsvc
	}
	if nk.ConfigMap != nil {
		cm := configuration.NewConfigMap()
		cm.K8sConfigMap = *nk.ConfigMap
		cm.ConvertK8sConfigMapToDB()
		cw.ConfigMap = &cm
	}
	if nk.Ingress != nil {
		ing := ingresses.NewIngress()
		ing.K8sIngress = *nk.Ingress
		err := ing.ConvertK8sIngressToDB()
		if err != nil {
			return cw, err
		}
		cw.Ingress = &ing
	}
	if nk.ServiceMonitor != nil {
		sm := servicemonitors.NewServiceMonitor()
		sm.K8sServiceMonitor = *nk.ServiceMonitor
		err := sm.ConvertK8sServiceMonitorToDB()
		if err != nil {
			return cw, err
		}
		cw.ServiceMonitor = &sm
	}
	if nk.Job != nil {
		j := jobs.NewJob()
		j.K8sJob = nk.Job
	}
	if nk.CronJob != nil {
		j := jobs.NewCronJob()
		j.K8sCronJob = nk.CronJob
	}
	return cw, nil
}

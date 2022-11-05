package chart_workload

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

type NativeK8s struct {
	*v1core.Service       `json:"service"`
	*v1core.ConfigMap     `json:"configMap"`
	*v1.Deployment        `json:"deployment"`
	*v1.StatefulSet       `json:"statefulSet"`
	*v1networking.Ingress `json:"ingress"`
}

func NewNativeK8s() NativeK8s {
	return NativeK8s{
		Service:     nil,
		ConfigMap:   nil,
		Deployment:  nil,
		StatefulSet: nil,
		Ingress:     nil,
	}
}

func (nk *NativeK8s) CreateChartWorkloadFromNativeK8s() (ChartWorkload, error) {
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
		sts := statefulset.NewStatefulSet()
		sts.K8sStatefulSet = *nk.StatefulSet
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
	return cw, nil
}

package zeus_core

import (
	"context"

	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v1networking "k8s.io/api/networking/v1"
)

type NamespaceWorkload struct {
	*v1.PodList               `json:"podList"`
	*v1.ServiceList           `json:"serviceList"`
	*v1networking.IngressList `json:"ingressList"`
	*v1apps.StatefulSetList   `json:"statefulSetList"`
	*v1apps.DeploymentList    `json:"deploymentList"`
	*v1.ConfigMapList         `json:"configMapList"`
}

func (k *K8Util) GetWorkloadAtNamespace(ctx context.Context, kns zeus_common_types.CloudCtxNs) (NamespaceWorkload, error) {
	k.SetContext(kns.Context)

	wrkLoad := NamespaceWorkload{}
	pods, err := k.GetPodsUsingCtxNs(ctx, kns, nil, nil)
	if err != nil {
		return wrkLoad, err
	}
	wrkLoad.PodList = pods

	svcs, err := k.GetServiceListWithKns(ctx, kns, nil)
	if err != nil {
		return wrkLoad, err
	}
	wrkLoad.ServiceList = svcs

	stss, err := k.GetStatefulSetList(ctx, kns, nil)
	if err != nil {
		return wrkLoad, err
	}
	wrkLoad.StatefulSetList = stss

	deps, err := k.GetDeploymentList(ctx, kns, nil)
	if err != nil {
		return wrkLoad, err
	}
	wrkLoad.DeploymentList = deps

	cfgMaps, err := k.GetConfigMapListWithKns(ctx, kns, nil)
	if err != nil {
		return wrkLoad, err
	}
	wrkLoad.ConfigMapList = cfgMaps
	return wrkLoad, err
}

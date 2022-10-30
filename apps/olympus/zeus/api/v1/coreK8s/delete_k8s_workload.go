package coreK8s

import (
	"context"

	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

func DeleteK8sWorkload(ctx context.Context, kns autok8s_core.KubeCtxNs, c read_charts.Chart) error {
	if c.Deployment != nil {
		// TODO
		_, err := K8util.DeleteDeployment(ctx, kns, &c.K8sDeployment)
		if err != nil {
			return err
		}
	}
	if c.Service != nil {
		// TODO
		err := K8util.DeleteServiceWithKns(ctx, kns, c.K8sService.Name, nil)
		if err != nil {
			return err
		}
	}
	if c.Ingress != nil {
		// TODO
		err := K8util.DeleteIngressWithKns(ctx, kns, c.K8sIngress.Name, nil)
		if err != nil {
			return err
		}
	}
	if c.ConfigMap != nil {
		// TODO
		err := K8util.DeleteConfigMapWithKns(ctx, kns, c.K8sConfigMap.Name, nil)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
	}
	return nil
}

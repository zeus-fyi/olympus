package delete_deploy

import (
	"context"

	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

func DeleteK8sWorkload(ctx context.Context, kns autok8s_core.KubeCtxNs, c read_charts.Chart) error {

	if c.Deployment != nil {
		err := core.K8Util.DeleteDeployment(ctx, kns, c.K8sDeployment.Name, nil)
		if err != nil {
			return err
		}
	}
	if c.Service != nil {
		err := core.K8Util.DeleteServiceWithKns(ctx, kns, c.K8sService.Name, nil)
		if err != nil {
			return err
		}
	}
	if c.ConfigMap != nil {
		err := core.K8Util.DeleteConfigMapWithKns(ctx, kns, c.K8sConfigMap.Name, nil)
		if err != nil {
			return err
		}
	}
	if c.Ingress != nil {
		err := core.K8Util.DeleteIngressWithKns(ctx, kns, c.K8sIngress.Name, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

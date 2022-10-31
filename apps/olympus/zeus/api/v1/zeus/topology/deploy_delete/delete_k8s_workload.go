package deploy_delete

import (
	"context"

	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

func DeleteK8sWorkload(ctx context.Context, kns autok8s_core.KubeCtxNs, c read_charts.Chart) error {

	// todo verify if it returns an error if not found, should ignore in those cases
	if c.Deployment != nil {
		// TODO
		err := core.K8Util.DeleteDeployment(ctx, kns, c.K8sDeployment.Name, nil)
		if err != nil {
			return err
		}
	}
	if c.Service != nil {
		// TODO
		err := core.K8Util.DeleteServiceWithKns(ctx, kns, c.K8sService.Name, nil)
		if err != nil {
			return err
		}
	}
	if c.Ingress != nil {
		// TODO
		err := core.K8Util.DeleteIngressWithKns(ctx, kns, c.K8sIngress.Name, nil)
		if err != nil {
			return err
		}
	}
	if c.ConfigMap != nil {
		// TODO
		err := core.K8Util.DeleteConfigMapWithKns(ctx, kns, c.K8sConfigMap.Name, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

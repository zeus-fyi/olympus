package delete_deploy

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

func DeleteK8sWorkload(ctx context.Context, kns autok8s_core.KubeCtxNs, nk chart_workload.NativeK8s) error {

	if nk.Deployment != nil {
		err := core.K8Util.DeleteDeployment(ctx, kns, nk.Deployment.Name, nil)
		if err != nil {
			return err
		}
	}
	if nk.Service != nil {
		err := core.K8Util.DeleteServiceWithKns(ctx, kns, nk.Service.Name, nil)
		if err != nil {
			return err
		}
	}
	if nk.ConfigMap != nil {
		err := core.K8Util.DeleteConfigMapWithKns(ctx, kns, nk.ConfigMap.Name, nil)
		if err != nil {
			return err
		}
	}
	if nk.Ingress != nil {
		err := core.K8Util.DeleteIngressWithKns(ctx, kns, nk.Ingress.Name, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

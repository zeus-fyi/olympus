package deploy

import (
	"context"

	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

func DeployChartPackage(ctx context.Context, kns autok8s_core.KubeCtxNs, c read_charts.Chart) error {

	_, nserr := core.K8Util.CreateNamespaceIfDoesNotExist(ctx, kns)
	if nserr != nil {
		return nserr
	}

	if c.Deployment != nil {
		// TODO
		_, err := core.K8Util.CreateDeploymentIfVersionLabelChangesOrDoesNotExist(ctx, kns, &c.K8sDeployment, nil)
		if err != nil {
			return err
		}
	}
	if c.Service != nil {
		// TODO
		_, err := core.K8Util.CreateServiceIfVersionLabelChangesOrDoesNotExist(ctx, kns, &c.K8sService, nil)
		if err != nil {
			return err
		}
	}
	if c.ConfigMap != nil {
		// TODO
		_, err := core.K8Util.CreateConfigMapIfVersionLabelChangesOrDoesNotExist(ctx, kns, &c.K8sConfigMap, nil)
		if err != nil {
			return err
		}
	}
	if c.Ingress != nil {
		// TODO
		_, err := core.K8Util.CreateIngressIfVersionLabelChangesOrDoesNotExist(ctx, kns, &c.K8sIngress, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

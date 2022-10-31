package coreK8s

import (
	"context"

	read_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/charts"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/zeus_pkg"
)

func DeployChartPackage(ctx context.Context, kns autok8s_core.KubeCtxNs, c read_charts.Chart) error {

	_, nserr := zeus_pkg.K8Util.CreateNamespaceIfDoesNotExist(ctx, kns)
	if nserr != nil {
		return nserr
	}

	if c.Deployment != nil {
		// TODO
		_, err := zeus_pkg.K8Util.CreateDeploymentIfVersionLabelChangesOrDoesNotExist(ctx, kns, &c.K8sDeployment, nil)
		if err != nil {
			return err
		}
	}
	if c.Service != nil {
		// TODO
		_, err := zeus_pkg.K8Util.CreateServiceIfVersionLabelChangesOrDoesNotExist(ctx, kns, &c.K8sService, nil)
		if err != nil {
			return err
		}
	}
	if c.ConfigMap != nil {
		// TODO
		_, err := zeus_pkg.K8Util.CreateConfigMapIfVersionLabelChangesOrDoesNotExist(ctx, kns, &c.K8sConfigMap, nil)
		if err != nil {
			return err
		}
	}
	if c.Ingress != nil {
		// TODO
		_, err := zeus_pkg.K8Util.CreateIngressIfVersionLabelChangesOrDoesNotExist(ctx, kns, &c.K8sIngress, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

package create_or_update_deploy

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DeployChartPackageComponent(ctx context.Context, kns autok8s_core.KubeCtxNs, nk chart_workload.NativeK8s, workloadKind string) error {
	switch workloadKind {
	case "Namespace":
		_, nserr := zeus.K8Util.CreateNamespaceIfDoesNotExist(ctx, kns)
		if nserr != nil {
			return nserr
		}
	case "Deployment":
		if nk.Deployment != nil {
			_, err := zeus.K8Util.CreateDeploymentIfVersionLabelChangesOrDoesNotExist(ctx, kns, nk.Deployment, nil)
			if err != nil {
				return err
			}
		}
	case "StatefulSet":
		if nk.StatefulSet != nil {
			_, err := zeus.K8Util.CreateStatefulSetIfVersionLabelChangesOrDoesNotExist(ctx, kns, nk.StatefulSet, nil)
			if err != nil {
				return err
			}
		}
	case "ConfigMap":
		if nk.ConfigMap != nil {
			_, err := zeus.K8Util.CreateConfigMapIfVersionLabelChangesOrDoesNotExist(ctx, kns, nk.ConfigMap, nil)
			if err != nil {
				return err
			}
		}
	case "Service":
		if nk.Service != nil {
			_, err := zeus.K8Util.CreateServiceIfVersionLabelChangesOrDoesNotExist(ctx, kns, nk.Service, nil)
			if err != nil {
				return err
			}
		}
	case "Ingress":
		if nk.Ingress != nil {
			_, err := zeus.K8Util.CreateIngressIfVersionLabelChangesOrDoesNotExist(ctx, kns, nk.Ingress, nil)
			if err != nil {
				return err
			}
		}
	default:
	}

	return nil
}

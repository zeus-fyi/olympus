package undeploy_topology

import (
	"context"
)

func (d *UndeployTopologyActivity) DeployIngress(ctx context.Context) error {
	if d.Ingress != nil {
		err := d.DeleteIngressWithKns(ctx, d.Kns, d.Ingress.Name, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

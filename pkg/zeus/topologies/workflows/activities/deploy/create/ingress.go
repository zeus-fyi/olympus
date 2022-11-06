package deploy_topology

import (
	"context"
)

func (d *DeployTopologyActivity) DeployIngress(ctx context.Context) error {
	if d.Ingress != nil {
		_, err := d.CreateIngressIfVersionLabelChangesOrDoesNotExist(ctx, d.Kns, d.Ingress, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

package deploy_topology

import (
	"context"
)

func (d *DeployTopologyActivity) DeployService(ctx context.Context) error {
	if d.Service != nil {
		_, err := d.CreateServiceIfVersionLabelChangesOrDoesNotExist(ctx, d.Kns, d.Service, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

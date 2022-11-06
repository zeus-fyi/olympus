package undeploy_topology

import (
	"context"
)

func (d *UndeployTopologyActivity) DeployService(ctx context.Context) error {
	if d.Service != nil {
		err := d.DeleteServiceWithKns(ctx, d.Kns, d.Service.Name, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

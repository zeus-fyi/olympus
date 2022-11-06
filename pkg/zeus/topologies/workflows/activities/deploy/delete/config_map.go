package undeploy_topology

import (
	"context"
)

func (d *UndeployTopologyActivity) DeployConfigMap(ctx context.Context) error {
	if d.ConfigMap != nil {
		err := d.DeleteConfigMapWithKns(ctx, d.Kns, d.ConfigMap.Name, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

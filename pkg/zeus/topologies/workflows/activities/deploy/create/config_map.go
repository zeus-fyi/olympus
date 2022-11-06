package deploy_topology

import (
	"context"
)

func (d *DeployTopologyActivity) DeployConfigMap(ctx context.Context) error {
	if d.ConfigMap != nil {
		_, err := d.CreateConfigMapIfVersionLabelChangesOrDoesNotExist(ctx, d.Kns, d.ConfigMap, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

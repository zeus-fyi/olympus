package deploy_topology

import (
	"context"
)

func (d *DeployTopologyActivity) DeployDeployment(ctx context.Context) error {
	if d.Deployment != nil {
		_, err := d.CreateDeploymentIfVersionLabelChangesOrDoesNotExist(ctx, d.Kns, d.Deployment, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

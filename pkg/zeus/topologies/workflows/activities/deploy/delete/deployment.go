package undeploy_topology

import (
	"context"
)

func (d *UndeployTopologyActivity) DeployDeployment(ctx context.Context) error {
	if d.Deployment != nil {
		err := d.DeleteDeployment(ctx, d.Kns, d.Deployment.Name, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

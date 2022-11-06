package deploy_topology

import (
	"context"
)

func (d *DeployTopologyActivity) CreateNamespace(ctx context.Context) error {
	_, nserr := d.CreateNamespaceIfDoesNotExist(ctx, d.Kns)
	if nserr != nil {
		return nserr
	}
	return nil
}

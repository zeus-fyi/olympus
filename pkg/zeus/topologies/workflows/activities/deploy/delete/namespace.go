package undeploy_topology

import (
	"context"
)

func (d *UndeployTopologyActivity) DestroyNamespace(ctx context.Context) error {
	nserr := d.DeleteNamespace(ctx, d.Kns)
	if nserr != nil {
		return nserr
	}
	return nil
}

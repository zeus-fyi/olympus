package destroy_deploy

import (
	"context"
)

func (d *UndeployTopologyActivity) DestroyNamespace(ctx context.Context) error {
	return d.postDestroyDeployTarget("namespace")
}

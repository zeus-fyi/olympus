package destroy_deploy

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DestroyNamespace(ctx context.Context) error {
	return d.postDestroyDeployTarget("namespace")
}

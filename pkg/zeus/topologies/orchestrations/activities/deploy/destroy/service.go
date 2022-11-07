package destroy_deploy

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DeployService(ctx context.Context) error {
	return d.postDestroyDeployTarget("service")
}

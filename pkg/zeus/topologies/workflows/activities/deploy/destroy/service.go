package destroy_deploy

import (
	"context"
)

func (d *UndeployTopologyActivity) DeployService(ctx context.Context) error {
	return d.postDestroyDeployTarget("service")
}

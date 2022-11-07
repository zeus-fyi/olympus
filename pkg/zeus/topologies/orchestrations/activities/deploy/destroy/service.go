package destroy_deploy

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DestroyDeployService(ctx context.Context) error {
	return d.postDestroyDeployTarget("service")
}

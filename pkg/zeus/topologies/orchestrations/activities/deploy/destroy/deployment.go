package destroy_deploy

import (
	"context"
)

func (d *DestroyDeployTopologyActivity) DestroyDeployDeployment(ctx context.Context) error {
	return d.postDestroyDeployTarget("deployment")
}

package deploy_topology

import (
	"context"

	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

//func WorkflowDeploy(ctx workflow.Context, params WorkflowParams) error {
//	var activities *YourActivityStruct
//	future := workflow.ExecuteActivity(ctx, activities.Activity, ctx)
//	var yourActivityResult YourActivityResult
//	if err := future.Get(ctx, &yourActivityResult); err != nil {
//		// ...
//	}
//	return nil
//}

type DeployParams struct {
	Kns        zeus_core.KubeCtxNs
	TopologyID int
	UserID     int
	OrgID      int
}

type DeployChart struct {
	K8Util zeus_core.K8Util

	DeployParams

	ActivityFieldOne string
	ActivityFieldTwo int
}

func (d *DeployChart) DeployWorkflow(ctx context.Context) error {

	return nil
}

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

type TopologyWorkflow struct {
}

func (t *TopologyWorkflow) TopologyWorkflow(ctx context.Context) error {

	//	var activities *YourActivityStruct
	//	future := workflow.ExecuteActivity(ctx, activities.Activity, ctx)
	//	var yourActivityResult YourActivityResult
	//	if err := future.Get(ctx, &yourActivityResult); err != nil {
	//		// ...
	//	}
	//	return nil
	return nil
}

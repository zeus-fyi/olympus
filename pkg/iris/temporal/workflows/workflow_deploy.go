package workflows

import (
	"context"

	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"go.temporal.io/sdk/workflow"
)

type WorkflowParams struct {
	Kns        zeus_core.KubeCtxNs
	TopologyID int
}

func WorkflowDeploy(ctx workflow.Context, params WorkflowParams) error {
	var activities *YourActivityStruct
	future := workflow.ExecuteActivity(ctx, activities.Activity, ctx)
	var yourActivityResult YourActivityResult
	if err := future.Get(ctx, &yourActivityResult); err != nil {
		// ...
	}
	return nil
}

type YourActivityResult struct {
}
type YourActivityStruct struct {
	ActivityFieldOne string
	ActivityFieldTwo int
}

func (a *YourActivityStruct) Activity(ctx context.Context) error {
	return nil
}

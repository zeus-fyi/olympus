package temporal_base

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

/*
 Workflows orchestrate the execution of Activities.
*/

type Workflow struct {
	Name       string
	WorkflowID string

	Params WorkflowParams
	Resp   WorkflowResponse

	WorkflowRun client.WorkflowRun
	StartOpts   client.StartWorkflowOptions
}

type WorkflowParams interface{}
type WorkflowResponse interface{}
type WorkflowDefinitionFn func(ctx workflow.Context, params WorkflowParams) (WorkflowResponse, error)

func (w *Workflow) SetAndReturnStartOpts(taskQueue string) client.StartWorkflowOptions {
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: taskQueue,
		ID:        w.WorkflowID,
	}
	w.StartOpts = workflowOptions
	return w.StartOpts
}

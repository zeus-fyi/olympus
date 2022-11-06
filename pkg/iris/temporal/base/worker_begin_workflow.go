package temporal_base

import "context"

func (w *Worker) BeginWorkflow(ctx context.Context, workFlow Workflow) error {
	err := w.Connect()
	if err != nil {
		return err
	}
	defer w.Close()
	// verify workFlow = YourWorkflowDefinition
	workflowRun, err := w.ExecuteWorkflow(ctx, workFlow.SetAndReturnStartOpts(w.TaskQueueName), workFlow, workFlow.Params)
	if err != nil {
		return err
	}
	workFlow.WorkflowRun = workflowRun
	return err
}

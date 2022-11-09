package temporal_base

import "context"

func (w *Worker) GetWorkflowResult(ctx context.Context, workFlow Workflow) error {
	err := w.ConnectTemporalClient()
	if err != nil {
		return err
	}
	defer w.Close()

	workflowID := workFlow.WorkflowID
	workflowRun := w.GetWorkflow(ctx, workflowID, workFlow.WorkflowRun.GetRunID())
	// verify workFlow = YourWorkflowDefinition
	err = workflowRun.Get(ctx, &workFlow.Resp)
	if err != nil {
		return err
	}
	workFlow.WorkflowRun = workflowRun
	return err
}

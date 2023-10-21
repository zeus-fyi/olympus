Kronos Readme

READ

use olympus group to get internal assignments
```go
func (k *KronosActivities) GetInternalAssignments(ctx context.Context) ([]artemis_orchestrations.OrchestrationJob, error) {
	ojs, err := artemis_orchestrations.SelectSystemOrchestrationsWithInstructionsByGroup(ctx, internalOrgID, olympus)
	if err != nil {
		return nil, err
	}
	return ojs, err
}
```

WRITE

ALERTS

see TestInsertAlertOrchestrators. 

groupName is the workflow group class, the inst type is the specific workflow

```go
	groupName := "HestiaPlatformServiceWorkflows"
	instType := "IrisRemoveAllOrgRoutesFromCacheWorkflow"

	orchName := fmt.Sprintf("%s-%s", groupName, instType)
	inst := Instructions{
		GroupName: groupName,
		Type:      instType,
		Alerts: AlertInstructions{
			Severity:  apollo_pagerduty.CRITICAL,
			Source:    "TEMPORAL_ALERTS",
			Component: "This is a workflow component",
			Message:   "A QuickNode services workflow is stuck",
		},
		Trigger: TriggerInstructions{
			AlertAfterTime:              time.Minute * 10,
			ResetAlertAfterTimeDuration: time.Minute * 10,
		},
	}
```

use UpsertAssignment to create a new assignment for kronos to watch
```go
func (k *KronosActivities) UpsertAssignment(ctx context.Context, oj artemis_orchestrations.OrchestrationJob) error {
	err := oj.UpsertOrchestrationWithInstructions(ctx)
	if err != nil {
		log.Err(err).Msg("UpsertAssignment: UpsertOrchestrationWithInstructions failed")
		return err
	}
	return nil
}
```

Example insertion to monitor workflow

EXECUTE_WORKFLOW

create a workflow id for tracking
```go
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaPlatformServiceWorkflows()
	wf := txWf.IrisDeleteOrgGroupRoutingTableWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisDeleteOrgGroupRoutingTableWorkflow")
		return err
	}
```

ENTRYPOINT
```go
    // or NewInternalActiveTemporalOrchestrationJobTemplate
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(internalOrgID, wfID, "HestiaPlatformServiceWorkflows", "IrisDeleteOrgGroupRoutingTableWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Warn("params", pr)
		logger.Warn("ou", pr.Ou)
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
```

EXIT
```go
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
    err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
    if err != nil {
        logger.Error("Failed to update and mark orchestration inactive", "Error", err)
    return err
```

START

use TestKronosHelixPattern, to start or restart Kronos

for internal platform alerts

use group: olympus
use type: alerts

see orchestrations jobs for database fns
datastores/postgres/apps/artemis/models/artemis_orchestrations/orchestrations.go
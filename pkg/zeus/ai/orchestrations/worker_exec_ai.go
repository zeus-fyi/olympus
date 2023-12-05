package ai_platform_service_orchestrations

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteRunAiWorkflowProcess(ctx context.Context, ou org_users.OrgUser, wfName string, delay time.Duration) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()

	wfID := CreateExecAiWfId(ou.OrgID, wfName)
	workflowOptions := client.StartWorkflowOptions{
		ID:        wfID,
		TaskQueue: z.TaskQueueName,
	}

	resp, _ := tc.DescribeWorkflowExecution(ctx, wfID, "")
	if resp != nil {
		// Check if the workflow is in a running state.
		if resp.WorkflowExecutionInfo.Status == enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			log.Warn().Msg("ExecuteIrisServerlessPodRestartWorkflow: workflow already running")
			return nil
		}
	}

	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.RunAiWorkflowProcess

	waitTill := time.Now().Add(delay)
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, wfID, ou, waitTill)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisServerlessPodRestartWorkflow")
		return err
	}
	return err
}

func (z *ZeusAiPlatformServicesWorker) EarlyStart(ctx context.Context, ou org_users.OrgUser, wfName string) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	wakeUpTime := time.Now()

	err := tc.SignalWorkflow(ctx, CreateExecAiWfId(ou.OrgID, wfName), "", SignalType, wakeUpTime)
	if err != nil {
		log.Err(err).Msg("IrisServicesWorker: EarlyStart")
		return err
	}
	return err
}

func CreateExecAiWfId(orgID int, wfName string) string {
	return fmt.Sprintf("%d-%s", orgID, wfName)
}

package iris_serverless

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

func (i *IrisServicesWorker) ExecuteIrisServerlessResyncWorkflow(ctx context.Context) error {
	tc := i.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: i.TaskQueueName,
	}
	txWf := NewIrisPlatformServiceWorkflows()
	wf := txWf.IrisServerlessResyncWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisServerlessResyncWorkflow")
		return err
	}
	return err
}

func (i *IrisServicesWorker) ExecuteIrisServerlessPodRestartWorkflow(ctx context.Context, orgID int, cctx zeus_common_types.CloudCtxNs, podName, serverlessTable, sessionID string, delay time.Duration) error {
	tc := i.ConnectTemporalClient()
	defer tc.Close()

	wfID := CreateServerlessRestartWfId(orgID, podName, serverlessTable, sessionID)
	workflowOptions := client.StartWorkflowOptions{
		ID:        wfID,
		TaskQueue: i.TaskQueueName,
	}

	resp, err := tc.DescribeWorkflowExecution(ctx, wfID, "")
	if err != nil {
		return err
	}
	// Check if the workflow is in a running state.
	if resp.WorkflowExecutionInfo.Status == enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
		log.Warn().Msg("ExecuteIrisServerlessPodRestartWorkflow: workflow already running")
		return nil
	}
	txWf := NewIrisPlatformServiceWorkflows()
	wf := txWf.IrisServerlessPodRestartWorkflow

	waitTill := time.Now().Add(delay)
	_, err = tc.ExecuteWorkflow(ctx, workflowOptions, wf, wfID, orgID, cctx, podName, serverlessTable, sessionID, waitTill)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisServerlessPodRestartWorkflow")
		return err
	}
	return err
}

func (i *IrisServicesWorker) EarlyStart(ctx context.Context, orgID int, podName, serverlessTable, sessionID string) error {
	tc := i.ConnectTemporalClient()
	defer tc.Close()
	wakeUpTime := time.Now()

	err := tc.SignalWorkflow(ctx, CreateServerlessRestartWfId(orgID, podName, serverlessTable, sessionID), "", SignalType, wakeUpTime)
	if err != nil {
		log.Err(err).Msg("IrisServicesWorker: EarlyStart")
		return err
	}
	return err
}

func CreateServerlessRestartWfId(orgID int, podName, serverlessTable, sessionID string) string {
	return fmt.Sprintf("%d-%s-%s-%s", orgID, podName, serverlessTable, sessionID)
}

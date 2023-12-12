package ai_platform_service_orchestrations

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteRunAiWorkflowProcess(ctx context.Context, ou org_users.OrgUser, params artemis_orchestrations.WorkflowExecParams) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	wfID := CreateExecAiWfId(params.WorkflowTemplate.WorkflowName)
	workflowOptions := client.StartWorkflowOptions{
		ID:        wfID,
		TaskQueue: z.TaskQueueName,
	}
	resp, _ := tc.DescribeWorkflowExecution(ctx, wfID, "")
	if resp != nil {
		// Check if the workflow is in a running state.
		if resp.WorkflowExecutionInfo.Status == enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			log.Warn().Msg("ExecuteRunAiWorkflowProcess: workflow already running")
			return nil
		}
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.RunAiWorkflowProcess
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, wfID, ou, params)
	if err != nil {
		log.Err(err).Msg("ExecuteRunAiWorkflowProcess")
		return err
	}
	return err
}

func (z *ZeusAiPlatformServicesWorker) EarlyStart(ctx context.Context, orchestrationName string) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	wakeUpTime := time.Now()
	err := tc.SignalWorkflow(ctx, orchestrationName, "", SignalType, wakeUpTime)
	if err != nil {
		log.Err(err).Msg("ZeusAiPlatformServicesWorker EarlyStartStop")
		return err
	}
	return err
}

func CreateExecAiWfId(wfName string) string {
	ud := uuid.New().String()
	return fmt.Sprintf("%s-%s-%s", wfName, strings.Split(ud, "-")[0], strings.Split(ud, "-")[1])
}

func (z *ZeusAiPlatformServicesWorker) ExecuteRunSearchIndexerWorkflowProcess(ctx context.Context, ou org_users.OrgUser, params artemis_orchestrations.WorkflowExecParams) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	wfID := CreateExecAiWfId(params.WorkflowTemplate.WorkflowName)
	workflowOptions := client.StartWorkflowOptions{
		ID:        wfID,
		TaskQueue: z.TaskQueueName,
	}
	resp, _ := tc.DescribeWorkflowExecution(ctx, wfID, "")
	if resp != nil {
		// Check if the workflow is in a running state.
		if resp.WorkflowExecutionInfo.Status == enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			log.Warn().Msg("ExecuteRunAiWorkflowProcess: workflow already running")
			return nil
		}
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.RunAiWorkflowProcess
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, wfID, ou, params)
	if err != nil {
		log.Err(err).Msg("ExecuteRunAiWorkflowProcess")
		return err
	}
	return err
}

func (z *ZeusAiPlatformServicesWorker) ExecuteAiSearchIndexerActionsWorkflow(ctx context.Context, ou org_users.OrgUser, params SearchIndexerActionsRequest) error {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("search-indexer-actions-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.AiSearchIndexerActionsWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, ou, params)
	if err != nil {
		log.Err(err).Msg("ExecuteAiSearchIndexerActionsWorkflow")
		return err
	}
	return nil
}

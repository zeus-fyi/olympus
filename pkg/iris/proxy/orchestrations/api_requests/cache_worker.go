package iris_api_requests

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (t *IrisApiRequestsWorker) ExecuteIrisCacheRefreshAllOrgRoutingTablesWorkflow(ctx context.Context) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue:          t.TaskQueueName,
		WorkflowRunTimeout: time.Hour,
	}
	txWf := NewIrisApiRequestsWorkflow()
	wf := txWf.CacheRefreshAllOrgRoutingTablesWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisCacheRefreshAllOrgRoutingTablesWorkflow")
		return err
	}
	return err
}

func (t *IrisApiRequestsWorker) ExecuteIrisCacheUpdateOrAddOrgRoutingTablesWorkflow(ctx context.Context, orgID int) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewIrisApiRequestsWorkflow()
	wf := txWf.CacheRefreshOrgRoutingTablesWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, orgID)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisCacheUpdateOrAddOrgRoutingTablesWorkflow")
		return err
	}
	return err
}

func (t *IrisApiRequestsWorker) ExecuteIrisCacheUpdateOrAddOrgGroupRoutingTableWorkflow(ctx context.Context, orgID int, groupName string) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewIrisApiRequestsWorkflow()
	wf := txWf.CacheRefreshOrgGroupTableWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, orgID, groupName)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisCacheUpdateOrAddOrgGroupRoutingTableWorkflow")
		return err
	}
	return err
}

func (t *IrisApiRequestsWorker) ExecuteIrisCacheDeleteOrgGroupRoutingTableWorkflow(ctx context.Context, orgID int, groupName string) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewIrisApiRequestsWorkflow()
	wf := txWf.DeleteRoutingGroupWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, orgID, groupName)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisCacheDeleteOrgGroupRoutingTableWorkflow")
		return err
	}
	return err
}

func (t *IrisApiRequestsWorker) ExecuteIrisCacheDeleteOrgRoutingGroupTablesWorkflow(ctx context.Context, orgID int) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewIrisApiRequestsWorkflow()
	wf := txWf.DeleteAllOrgRoutingGroupsWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, orgID)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisCacheDeleteAllOrRoutingGroupTablesWorkflow")
		return err
	}
	return err
}

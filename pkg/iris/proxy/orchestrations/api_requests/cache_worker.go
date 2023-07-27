package iris_api_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (t *IrisApiRequestsWorker) ExecuteIrisCacheUpdateOrAddOrgRoutingTablesWorkflow(ctx context.Context, orgID int) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	txWf := NewIrisApiRequestsWorkflow()
	wf := txWf.CacheUpdateOrAddOrgRoutingTablesWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, orgID)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisCacheUpdateOrAddOrgRoutingTablesWorkflow")
		return err
	}
	return err
}

func (t *IrisApiRequestsWorker) ExecuteIrisCacheRefreshAllOrgRoutingTablesWorkflow(ctx context.Context) error {
	tc := t.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
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

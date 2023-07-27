package iris_api_requests

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

func (i *IrisApiRequestsWorkflow) CacheUpdateOrAddOrgRoutingTablesWorkflow(ctx workflow.Context, orgID int) error {
	err := i.UpdateOrgRoutingTables(context.Background(), orgID)
	if err != nil {
		return err
	}
	return nil
}

func (i *IrisApiRequestsWorkflow) CacheRefreshAllOrgRoutingTablesWorkflow(ctx workflow.Context) error {
	err := i.RefreshAllOrgRoutingTables(context.Background())
	if err != nil {
		return err
	}
	return nil
}

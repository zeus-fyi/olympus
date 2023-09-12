package iris_api_requests

import (
	"time"

	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (i *IrisApiRequestsWorkflow) CacheRefreshAllOrgRoutingTablesWorkflow(ctx workflow.Context) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    15 * time.Minute,
			BackoffCoefficient: 2,
			MaximumAttempts:    20,
		},
	}
	getRoutingTablesCtx := workflow.WithActivityOptions(ctx, ao)
	var ogr iris_models.OrgRoutesGroup
	err := workflow.ExecuteActivity(getRoutingTablesCtx, i.SelectAllRoutingTables).Get(getRoutingTablesCtx, &ogr)
	if err != nil {
		log.Error("CacheRefreshAllOrgRoutingTablesWorkflow: Failed to SelectAllRoutingTables", "Error", err)
		return err
	}

	for orgID, og := range ogr.Map {
		for rgName, routes := range og {
			if rgName == "unused" {
				continue
			}
			addOrUpdateOrgRoutingTableCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(addOrUpdateOrgRoutingTableCtx, i.UpdateOrgRoutingTable, orgID, rgName, routes).Get(addOrUpdateOrgRoutingTableCtx, nil)
			if err != nil {
				log.Error("CacheRefreshAllOrgRoutingTablesWorkflow: Failed to UpdateOrgRoutingTable", "Error", err)
				return err
			}
		}
	}
	return nil
}

func (i *IrisApiRequestsWorkflow) CacheRefreshOrgRoutingTablesWorkflow(ctx workflow.Context, orgID int) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Minute,
			BackoffCoefficient: 2,
		},
	}
	getRoutingTablesCtx := workflow.WithActivityOptions(ctx, ao)
	var ogr map[string][]iris_models.RouteInfo
	err := workflow.ExecuteActivity(getRoutingTablesCtx, i.SelectSingleOrgGroupsRoutingTables, orgID).Get(getRoutingTablesCtx, &ogr)
	if err != nil {
		log.Error("CacheRefreshOrgRoutingTablesWorkflow: Failed to SelectAllOrgGroupsRoutingTables", "Error", err)
		return err
	}
	for rgName, routes := range ogr {
		if rgName == "unused" {
			continue
		}
		addOrUpdateOrgRoutingTableCtx := workflow.WithActivityOptions(ctx, ao)
		err = workflow.ExecuteActivity(addOrUpdateOrgRoutingTableCtx, i.UpdateOrgRoutingTable, orgID, rgName, routes).Get(addOrUpdateOrgRoutingTableCtx, nil)
		if err != nil {
			log.Error("CacheRefreshOrgRoutingTablesWorkflow: Failed to UpdateOrgRoutingTable", "Error", err)
			return err
		}
	}
	return nil
}

func (i *IrisApiRequestsWorkflow) CacheRefreshOrgGroupTableWorkflow(ctx workflow.Context, orgID int, groupName string) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 30, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    3 * time.Minute,
			BackoffCoefficient: 2,
		},
	}
	getRoutingTablesCtx := workflow.WithActivityOptions(ctx, ao)
	var ogr iris_models.OrgRoutesGroup
	err := workflow.ExecuteActivity(getRoutingTablesCtx, i.SelectOrgGroupRoutingTable, orgID, groupName).Get(getRoutingTablesCtx, &ogr)
	if err != nil {
		log.Error("CacheRefreshOrgGroupTableWorkflow: Failed to SelectOrgGroupRoutingTable", "Error", err)
		return err
	}
	for _, og := range ogr.Map {
		for rgName, routes := range og {
			if rgName == "unused" {
				continue
			}
			addOrUpdateOrgRoutingTableCtx := workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(addOrUpdateOrgRoutingTableCtx, i.UpdateOrgRoutingTable, orgID, rgName, routes).Get(addOrUpdateOrgRoutingTableCtx, nil)
			if err != nil {
				log.Error("CacheRefreshOrgGroupTableWorkflow: Failed to UpdateOrgRoutingTable", "Error", err)
				return err
			}
		}
	}
	return nil
}

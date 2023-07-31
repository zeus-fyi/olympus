package iris_api_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
)

func (i *IrisApiRequestsActivities) UpdateOrgRoutingTable(ctx context.Context, orgID int, rgName string, routes []iris_models.RouteInfo) error {
	err := iris_redis.IrisRedis.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
	if err != nil {
		log.Error().Int("orgID", orgID).Str("routeGroup", rgName).Err(err).Msg("UpdateOrgRoutingTable: Failed to update routing tables for org")
		return err
	}
	return nil
}

func (i *IrisApiRequestsActivities) DeleteOrgRoutingTable(ctx context.Context, orgID int, rgName string) error {
	err := iris_redis.IrisRedis.DeleteOrgRoutingGroup(context.Background(), orgID, rgName)
	if err != nil {
		log.Error().Int("orgID", orgID).Str("routeGroup", rgName).Err(err).Msg("UpdateOrgRoutingTable: Failed to update routing tables for org")
		return err
	}
	return nil
}

func (i *IrisApiRequestsActivities) SelectOrgGroupRoutingTable(ctx context.Context, orgID int, groupName string) (iris_models.OrgRoutesGroup, error) {
	ot, err := iris_models.SelectOrgRoutesByOrgAndGroupName(ctx, orgID, groupName)
	if err != nil {
		return ot, err
	}
	return ot, nil
}

func (i *IrisApiRequestsActivities) SelectAllOrgGroupsRoutingTables(ctx context.Context, orgID int) (iris_models.OrgRoutesGroup, error) {
	ot, err := iris_models.SelectAllOrgRoutesByOrg(ctx, orgID)
	if err != nil {
		return ot, err
	}
	return ot, nil
}

func (i *IrisApiRequestsActivities) SelectAllRoutingTables(ctx context.Context) (iris_models.OrgRoutesGroup, error) {
	ot, err := iris_models.SelectAllOrgRoutes(ctx)
	if err != nil {
		return ot, err
	}
	return ot, nil
}

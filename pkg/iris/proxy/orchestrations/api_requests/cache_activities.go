package iris_api_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
)

func (i *IrisApiRequestsActivities) UpdateOrgRoutingTables(ctx context.Context, orgID int) error {
	err := iris_redis.IrisRedis.InitRoutingTablesForOrg(context.Background(), orgID)
	if err != nil {
		log.Error().Int("orgID", orgID).Err(err).Msg("UpdateOrgRoutingTables: Failed to update routing tables for org")
		return err
	}
	return nil
}

func (i *IrisApiRequestsActivities) RefreshAllOrgRoutingTables(ctx context.Context) error {
	err := iris_redis.IrisRedis.InitRoutingTables(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("RefreshAllOrgRoutingTables: Failed to refresh routing tables")
		return err
	}
	return nil
}

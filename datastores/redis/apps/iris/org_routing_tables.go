package iris_redis

import (
	"context"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (m *IrisCache) AddOrUpdateOrgRoutingGroup(ctx context.Context, orgID int, rgName string, routes []iris_models.RouteInfo) error {
	// Start a new transaction
	pipe := m.Writer.TxPipeline()

	// Define the route group tag
	rgTag := getOrgRouteKey(orgID, rgName)

	// Remove the old route group key if it exists
	pipe.Del(ctx, rgTag)

	// Add each route to the list
	for _, routeInfo := range routes {
		// Define unique tags for route and referers
		routeTag := getOrgRouteKey(orgID, routeInfo.RoutePath)
		refererTag := routeTag + ":referers"

		// Remove the old keys if they exist
		pipe.Del(ctx, routeTag, refererTag)

		// Add the route
		pipe.RPush(ctx, rgTag, routeInfo.RoutePath)

		// Add the referers to a set
		if len(routeInfo.Referrers) > 0 {
			// Convert []string to []interface{}
			referersInterface := make([]interface{}, len(routeInfo.Referrers))
			for i, v := range routeInfo.Referrers {
				referersInterface[i] = v
			}
			pipe.SAdd(ctx, refererTag, referersInterface...)
		}
	}

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("AddOrUpdateOrgRoutingGroup: %s", rgTag)
		return err
	}
	return nil
}

func (m *IrisCache) DeleteOrgRoutingGroup(ctx context.Context, orgID int, rgName string) error {
	// Generate the key
	tag := getOrgRouteKey(orgID, rgName)

	// Start a new transaction
	pipe := m.Writer.TxPipeline()

	// Remove the old key if it exists
	pipe.Del(ctx, tag)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("DeleteOrgRoutingGroup: %s", tag)
		return err
	}
	return nil
}

func (m *IrisCache) initRoutingTables(ctx context.Context) error {
	ot, err := iris_models.SelectAllOrgRoutes(ctx)
	if err != nil {
		return err
	}
	for orgID, og := range ot.Map {
		for rgName, routes := range og {
			err = m.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
			if err != nil {
				log.Err(err).Msg("InitRoutingTables")
				return err
			}
		}
	}
	return nil
}

func (m *IrisCache) initRoutingTablesForOrg(ctx context.Context, orgID int) error {
	ot, err := iris_models.SelectAllOrgRoutesByOrg(ctx, orgID)
	if err != nil {
		return err
	}
	for rgName, routes := range ot {
		err = m.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
		if err != nil {
			log.Err(err).Msg("InitRoutingTablesForOrg")
			return err
		}
	}
	return nil
}

func (m *IrisCache) refreshRoutingTablesForOrgGroup(ctx context.Context, orgID int, groupName string) error {
	ot, err := iris_models.SelectOrgRoutesByOrgAndGroupName(ctx, orgID, groupName)
	if err != nil {
		return err
	}
	for _, og := range ot.Map {
		for rgName, routes := range og {
			err = m.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
			if err != nil {
				log.Err(err).Msg("InitRoutingTablesForOrgGroup")
				return err
			}
		}
	}
	return nil
}
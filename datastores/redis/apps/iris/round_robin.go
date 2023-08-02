package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func orgRouteTag(orgID int, rgName string) string {
	return fmt.Sprintf("%d-%s", orgID, rgName)
}

func orgMonthlyUsageTag(orgID int, month string) string {
	return fmt.Sprintf("%d-%s-total-zu-usage", orgID, month)
}

func orgRateLimitTag(orgID int) string {
	return fmt.Sprintf("%d:%d", orgID, time.Now().Unix())
}

func (m *IrisCache) GetNextRoute(ctx context.Context, orgID int, rgName string, meter *iris_usage_meters.PayloadSizeMeter) (iris_models.RouteInfo, error) {
	// Generate the rate limiter key with the Unix timestamp
	rateLimiterKey := orgRateLimitTag(orgID)
	orgRequests := fmt.Sprintf("%d-total-zu-usage", orgID)
	if meter != nil {
		orgRequests = orgMonthlyUsageTag(orgID, meter.Month)
	}

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	if meter != nil {
		// Increment the payload size meter
		_ = pipe.IncrByFloat(ctx, orgRequests, meter.ZeusRequestComputeUnitsConsumed())
		// Increment the rate limiter key
		_ = pipe.IncrByFloat(ctx, rateLimiterKey, meter.ZeusRequestComputeUnitsConsumed())
	} else {
		_ = pipe.Incr(ctx, orgRequests)
		// Increment the rate limiter key
		_ = pipe.Incr(ctx, rateLimiterKey)
	}

	// Set the key to expire after 2 seconds
	pipe.Expire(ctx, rateLimiterKey, 10*time.Second)

	// Generate the route key
	routeKey := orgRouteTag(orgID, rgName)

	// Check if the key exists
	existsCmd := pipe.Exists(ctx, routeKey)

	// Pop the endpoint from the head of the list and push it back to the tail, ensuring round-robin rotation
	endpointCmd := pipe.LMove(ctx, routeKey, routeKey, "LEFT", "RIGHT")

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		// If there's an error, log it and return it
		fmt.Printf("error during pipeline execution: %s\n", err.Error())
		return iris_models.RouteInfo{}, err
	}

	// Check whether the key exists
	if existsCmd.Val() <= 0 {
		log.Err(err).Msgf("key doesn't exist: %s", routeKey)
		return iris_models.RouteInfo{}, fmt.Errorf("key doesn't exist: %s", routeKey)
	}

	// Get referers from Redis
	refererKey := orgRouteTag(orgID, endpointCmd.Val()) + ":referers"
	referers, err := m.Reader.SMembers(ctx, refererKey).Result()
	if err != nil {
		log.Err(err).Msgf("error getting referers for route: %s", refererKey)
		return iris_models.RouteInfo{}, fmt.Errorf("error getting referers for route: %s", refererKey)
	}

	// Return the RouteInfo
	return iris_models.RouteInfo{
		RoutePath: endpointCmd.Val(),
		Referers:  referers,
	}, nil
}

func (m *IrisCache) IncrementResponseUsageRateMeter(ctx context.Context, orgID int, meter *iris_usage_meters.PayloadSizeMeter) error {
	// Generate the rate limiter key with the Unix timestamp
	rateLimiterKey := orgRateLimitTag(orgID)
	orgRequests := fmt.Sprintf("%d-total-zu-usage", orgID)
	if meter != nil {
		orgRequests = orgMonthlyUsageTag(orgID, meter.Month)
	}

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()
	if meter != nil {
		// Increment the payload size meter
		_ = pipe.IncrByFloat(ctx, orgRequests, meter.ZeusResponseComputeUnitsConsumed())
		// Increment the rate limiter key
		_ = pipe.IncrByFloat(ctx, rateLimiterKey, meter.ZeusResponseComputeUnitsConsumed())
	} else {
		_ = pipe.Incr(ctx, orgRequests)
		// Increment the rate limiter key
		_ = pipe.Incr(ctx, rateLimiterKey)
	}
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("AddOrUpdateOrgRoutingGroup")
		return err
	}
	return nil
}

func (m *IrisCache) AddOrUpdateOrgRoutingGroup(ctx context.Context, orgID int, rgName string, routes []iris_models.RouteInfo) error {
	// Start a new transaction
	pipe := m.Writer.TxPipeline()

	// Define the route group tag
	rgTag := orgRouteTag(orgID, rgName)

	// Remove the old route group key if it exists
	pipe.Del(ctx, rgTag)

	// Add each route to the list
	for _, routeInfo := range routes {
		// Define unique tags for route and referers
		routeTag := orgRouteTag(orgID, routeInfo.RoutePath)
		refererTag := routeTag + ":referers"

		// Remove the old keys if they exist
		pipe.Del(ctx, routeTag, refererTag)

		// Add the route
		pipe.RPush(ctx, rgTag, routeInfo.RoutePath)

		// Add the referers to a set
		if len(routeInfo.Referers) > 0 {
			// Convert []string to []interface{}
			referersInterface := make([]interface{}, len(routeInfo.Referers))
			for i, v := range routeInfo.Referers {
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
	tag := orgRouteTag(orgID, rgName)

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

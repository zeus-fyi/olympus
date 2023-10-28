package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_usage_meters "github.com/zeus-fyi/olympus/pkg/iris/proxy/usage_meters"
)

func (m *IrisCache) GetNextRoute(ctx context.Context, orgID int, rgName string, meter *iris_usage_meters.PayloadSizeMeter) (iris_models.RouteInfo, error) {
	// Generate the rate limiter key with the Unix timestamp
	rateLimiterKey := getOrgRateLimitKey(orgID)
	orgRequests := getOrgTotalRequestsKey(orgID)
	if meter != nil {
		orgRequests = getOrgMonthlyUsageKey(orgID, meter.Month)
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

	// Generate the route key
	routeKey := getOrgRouteKey(orgID, rgName)

	// Check if the key exists
	existsCmd := pipe.Exists(ctx, routeKey)

	// Pop the endpoint from the head of the list and push it back to the tail, ensuring round-robin rotation
	endpointCmd := pipe.LMove(ctx, routeKey, routeKey, "LEFT", "RIGHT")

	// Set the key to expire after 2 seconds
	pipe.Expire(ctx, rateLimiterKey, 2*time.Second)

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
	refererKey := getOrgRouteKey(orgID, endpointCmd.Val()) + ":referers"
	referers, err := m.Reader.SMembers(ctx, refererKey).Result()
	if err != nil {
		log.Err(err).Msgf("error getting referers for route: %s", refererKey)
		return iris_models.RouteInfo{}, fmt.Errorf("error getting referers for route: %s", refererKey)
	}

	// Return the RouteInfo
	return iris_models.RouteInfo{
		RoutePath: endpointCmd.Val(),
		Referrers: referers,
	}, nil
}

func (m *IrisCache) IncrementResponseUsageRateMeter(ctx context.Context, orgID int, meter *iris_usage_meters.PayloadSizeMeter) error {
	// Generate the rate limiter key with the Unix timestamp
	if meter == nil {
		log.Warn().Msg("IncrementResponseUsageRateMeter: meter is nil")
		return nil
	}

	rateLimiterKey := getOrgRateLimitKey(orgID)
	orgRequests := getOrgTotalRequestsKey(orgID)
	if meter != nil {
		orgRequests = getOrgMonthlyUsageKey(orgID, meter.Month)
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

	pipe.Expire(ctx, rateLimiterKey, 2*time.Second)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("AddOrUpdateOrgRoutingGroup")
		return err
	}
	return nil
}

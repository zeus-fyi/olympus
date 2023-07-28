package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func orgRouteTag(orgID int, rgName string) string {
	return fmt.Sprintf("%d-%s", orgID, rgName)
}

var LocalRateLimiterCache = cache.New(1*time.Second, 2*time.Second)

func (m *IrisCache) GetNextRoute(ctx context.Context, orgID int, rgName string) (string, error) {
	// Generate the rate limiter key with the Unix timestamp
	rateLimiterKey := fmt.Sprintf("%d:%d", orgID, time.Now().Unix())
	_, found := LocalRateLimiterCache.Get(rateLimiterKey)
	if found {
		err := LocalRateLimiterCache.Increment(rateLimiterKey, 1)
		if err != nil {
			log.Err(err).Msg("LocalRateLimiterCache: GetNextRoute")
		}
	} else {
		LocalRateLimiterCache.Set(rateLimiterKey, 1, 1*time.Second)
	}

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	// Increment the rate limiter key
	rateCmd := pipe.Incr(ctx, rateLimiterKey)

	// Set the key to expire after 3 seconds
	pipe.Expire(ctx, rateLimiterKey, 3*time.Second)

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
		return "", err
	}

	// Check if rate limit has been exceeded
	if rateCmd.Val() > 20 {
		return "", fmt.Errorf("rate limit exceeded for orgID: %d", orgID)
	}

	// Check whether the key exists
	if existsCmd.Val() <= 0 {
		log.Err(err).Msgf("key doesn't exist: %s", routeKey)
		return "", fmt.Errorf("key doesn't exist: %s", routeKey)
	}

	// Return the endpoint
	return endpointCmd.Val(), nil
}

func (m *IrisCache) AddOrUpdateOrgRoutingGroup(ctx context.Context, orgID int, rgName string, routes []string) error {
	// Generate the key
	tag := orgRouteTag(orgID, rgName)

	// Start a new transaction
	pipe := m.Writer.TxPipeline()

	// Remove the old key if it exists
	pipe.Del(ctx, tag)

	// Add each route to the list
	for _, route := range routes {
		pipe.RPush(ctx, tag, route)
	}

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("AddOrUpdateOrgRoutingGroup: %s", tag)
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
	for _, og := range ot.Map {
		for rgName, routes := range og {
			err = m.AddOrUpdateOrgRoutingGroup(context.Background(), orgID, rgName, routes)
			if err != nil {
				log.Err(err).Msg("InitRoutingTablesForOrg")
				return err
			}
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

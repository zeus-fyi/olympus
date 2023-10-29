package iris_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

const (
	ServerlessAnvilTable = "anvil"

	// ServerlessSessionMaxRunTime adds one minute of margin to release the route
	ServerlessSessionMaxRunTime = 10 * time.Minute
	ServerlessMaxRunTime        = 11 * time.Minute
)

func (m *IrisCache) AddOrUpdateServerlessRoutingTable(ctx context.Context, serverlessRoutesTable string, routes []iris_models.RouteInfo) error {
	if len(routes) <= 0 {
		return m.RefreshServerlessRoutingTable(ctx, serverlessRoutesTable)
	}

	serviceAvailabilityTable := getGlobalServerlessAvailabilityTableKey(serverlessRoutesTable)
	tn := time.Now().Unix()
	pipe := m.Writer.TxPipeline()

	for _, r := range routes {
		redisSet := redis.Z{
			Score:  float64(tn),
			Member: r.RoutePath,
		}
		pipe.ZAddNX(ctx, serviceAvailabilityTable, redisSet)
	}

	minElemsCmd := pipe.ZRangeByScoreWithScores(ctx, serviceAvailabilityTable, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   fmt.Sprintf("%d", tn),
		Count: 100, // Get only the first route
	})
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return err
	}

	pipe = m.Writer.TxPipeline()
	members, err := minElemsCmd.Result()
	if err != nil {
		log.Err(err).Msgf("GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing")
		return err
	}
	serverlessReadyRoutes := getGlobalServerlessTableKey(serverlessRoutesTable)
	pipe = m.Writer.TxPipeline() // Starting a new pipeline for this batch of operations
	for _, member := range members {
		routePath, eok := member.Member.(string)
		if !eok {
			log.Err(fmt.Errorf("failed to convert member to string")).Msgf("Member: %v", routePath)
			continue
		}
		pipe.SAdd(ctx, serverlessReadyRoutes, routePath)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return err
	}
	return nil
}

func (m *IrisCache) RefreshServerlessRoutingTable(ctx context.Context, serverlessRoutesTable string) error {
	serviceAvailabilityTable := getGlobalServerlessAvailabilityTableKey(serverlessRoutesTable)
	pipe := m.Reader.TxPipeline()

	tn := time.Now().Unix()
	minElemsCmd := pipe.ZRangeByScoreWithScores(ctx, serviceAvailabilityTable, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   fmt.Sprintf("%d", tn),
		Count: 100, // Get only the first route
	})
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return err
	}

	pipe = m.Writer.TxPipeline()
	serverlessReadyRoutes := getGlobalServerlessTableKey(serverlessRoutesTable)
	members, err := minElemsCmd.Result()
	for _, member := range members {
		routePath, eok := member.Member.(string)
		if !eok {
			log.Err(fmt.Errorf("failed to convert member to string")).Msgf("Member: %v", routePath)
			continue
		}
		pipe.SAdd(ctx, serverlessReadyRoutes, routePath)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return err
	}
	return nil
}

// GetNextServerlessRoute TODO: needs to session count limit before here
func (m *IrisCache) GetNextServerlessRoute(ctx context.Context, orgID int, sessionID, serverlessRoutesTable string) (string, error) {
	serviceTable := getGlobalServerlessTableKey(serverlessRoutesTable)
	pipe := m.Writer.TxPipeline()

	// Pop the first item from the list
	routeCmd := pipe.SPop(ctx, serviceTable)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return "", err
	}
	path, err := routeCmd.Result()
	if err != nil {
		log.Err(err)
		return "", err
	}
	if len(path) <= 0 {
		return "", fmt.Errorf("GetNextServerlessRoute: failed to find any available routes")
	}
	pipe = m.Writer.TxPipeline()
	tn := time.Now().Add(ServerlessMaxRunTime).Unix()
	redisSet := redis.Z{
		Score:  float64(tn),
		Member: path,
	}
	// XX: Only update elements that already exist. Don't add new elements
	pipe.ZAddXX(ctx, getGlobalServerlessAvailabilityTableKey(serverlessRoutesTable), redisSet)

	// NX: Set key to hold string value if key does not exist.
	// this is used to k->v lookup for org-sessionID -> route
	pipe.SetNX(ctx, getOrgSessionIDKey(orgID, sessionID), path, ServerlessSessionMaxRunTime)

	// Add the session to the rate limit table
	sessionRateLimit := getOrgActiveServerlessCountKey(orgID, serverlessRoutesTable)
	tnUser := time.Now().Add(ServerlessSessionMaxRunTime).Unix()
	redisSetUser := redis.Z{
		Score:  float64(tnUser),
		Member: path,
	}
	// NX: Add to sorted set if key doesn't exist
	// Add the session to the rate limit table
	pipe.ZAddNX(ctx, sessionRateLimit, redisSetUser)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msg("GetNextServerlessRoute: failed to update availability table")
		return "", err
	}

	return path, nil
}

func (m *IrisCache) ReleaseServerlessRoute(ctx context.Context, orgID int, sessionID, serverlessRoutesTable string) error {
	serviceAvailabilityTable := getGlobalServerlessAvailabilityTableKey(serverlessRoutesTable)
	pipe := m.Reader.TxPipeline()

	tn := time.Now().Unix()
	minElemsCmd := pipe.ZRangeByScoreWithScores(ctx, serviceAvailabilityTable, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   fmt.Sprintf("%d", tn),
		Count: 100, // Get only the first route
	})
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return err
	}

	pipe = m.Writer.TxPipeline()
	serverlessReadyRoutes := getGlobalServerlessTableKey(serverlessRoutesTable)
	members, err := minElemsCmd.Result()
	for _, member := range members {
		routePath, eok := member.Member.(string)
		if !eok {
			log.Err(fmt.Errorf("failed to convert member to string")).Msgf("Member: %v", routePath)
			continue
		}
		pipe.SAdd(ctx, serverlessReadyRoutes, routePath)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return err
	}
	return nil
}

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

	TimeMarginBuffer = time.Minute

	// ServerlessSessionMaxRunTime adds one minute of margin to release the route
	ServerlessSessionMaxRunTime = 10 * time.Minute
	ServerlessMaxRunTime        = 10*time.Minute + TimeMarginBuffer
)

func (m *IrisCache) AddRoutesToServerlessRoutingTable(ctx context.Context, serverlessRoutesTable string, routes []iris_models.RouteInfo) error {
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
		// NX: Set key to hold string value if key does not exist.
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
		if err == redis.Nil {
			return "", fmt.Errorf("GetNextServerlessRoute: failed to find any available routes")
		}
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

	// expires the session rate limit table, adds 2x margin on ttl
	pipe.Expire(ctx, sessionRateLimit, ServerlessMaxRunTime*2)

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msg("GetNextServerlessRoute: failed to update availability table")
		return "", err
	}

	return path, nil
}

func (m *IrisCache) ReleaseServerlessRoute(ctx context.Context, orgID int, sessionID, serverlessRoutesTable string) error {
	pipe := m.Reader.TxPipeline()

	// Gets and removes session-to-route cache key
	orgSessionToRouteKey := getOrgSessionIDKey(orgID, sessionID)
	getSessionRoute := pipe.Get(ctx, orgSessionToRouteKey)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return err
	}
	// Get the value from the result of the Get command
	path, err := getSessionRoute.Result()
	if err != nil {
		log.Err(err).Msg("ReleaseServerlessRoute: failed to get session route")
		return err
	}
	pipe = m.Writer.TxPipeline()
	// deletes session -> route cache key
	pipe.Del(ctx, orgSessionToRouteKey)

	// Resets the timing of the route with some margin
	if len(path) > 0 {
		redisSet := redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: path,
		}
		// this is the user specific availability table
		sessionRateLimit := getOrgActiveServerlessCountKey(orgID, serverlessRoutesTable)
		pipe.ZRem(ctx, sessionRateLimit, path)

		// this is the global availability table
		// XX: Update elements that already exist. Don't add new elements
		serviceAvailabilityTable := getGlobalServerlessAvailabilityTableKey(serverlessRoutesTable)
		pipe.ZAddXX(ctx, serviceAvailabilityTable, redisSet)
		// this adds the route back to the set of available routes

		serverlessReadyRoutes := getGlobalServerlessTableKey(serverlessRoutesTable)
		pipe.SAdd(ctx, serverlessReadyRoutes, path)
	}

	// Removes expired sessions
	tn := time.Now().Unix()
	sessionRateLimit := getOrgActiveServerlessCountKey(orgID, serverlessRoutesTable)
	pipe.ZRemRangeByScore(ctx, sessionRateLimit, "0", fmt.Sprintf("%d", tn))

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return err
	}
	return nil
}

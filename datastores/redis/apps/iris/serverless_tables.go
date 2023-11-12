package iris_redis

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

const (
	ServerlessAnvilTable = "anvil"

	MaxActiveServerlessSessions = 5
	TimeMarginBuffer            = 2 * time.Minute

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
		if r.RoutePath == "" {
			continue
		}
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
		log.Err(err).Msg("AddRoutesToServerlessRoutingTable: failed to add routes to availability table")
		return err
	}
	members, err := minElemsCmd.Result()
	if err != nil {
		log.Err(err).Msgf("AddRoutesToServerlessRoutingTable")
		return err
	}
	if len(members) <= 0 {
		return nil
	}
	serverlessReadyRoutes := getGlobalServerlessTableKey(serverlessRoutesTable)
	pipe = m.Writer.TxPipeline() // Starting a new pipeline for this batch of operations
	for _, member := range members {
		routePath, eok := member.Member.(string)
		if !eok {
			log.Err(fmt.Errorf("failed to convert member to string")).Msgf("Member: %v", routePath)
			continue
		}
		if len(routePath) > 0 {
			pipe.SAdd(ctx, serverlessReadyRoutes, routePath)
		}
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msg("AddRoutesToServerlessRoutingTable: failed to add routes to availability table")
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
		Count: 100, // Get only the first routes
	})
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msg("RefreshServerlessRoutingTable: failed to get routes from availability table")
		return err
	}

	members, err := minElemsCmd.Result()
	if err != nil {
		log.Err(err).Msgf("RefreshServerlessRoutingTable")
		return err
	}
	if len(members) <= 0 {
		return nil
	}
	log.Info().Msgf("RefreshServerlessRoutingTable: found %d routes", len(members))
	serverlessReadyRoutes := getGlobalServerlessTableKey(serverlessRoutesTable)
	pipe = m.Writer.TxPipeline()
	for _, member := range members {
		routePath, eok := member.Member.(string)
		if !eok {
			log.Err(fmt.Errorf("RefreshServerlessRoutingTable: failed to convert member to string")).Msgf("Member: %v", routePath)
			continue
		}
		if len(routePath) <= 0 {
			continue
		}
		pipe.SAdd(ctx, serverlessReadyRoutes, routePath)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msg("RefreshServerlessRoutingTable: failed to refresh routes")
		return err
	}
	return nil
}

func (m *IrisCache) CheckServerlessSessionRateLimit(ctx context.Context, orgID int, sessionID, serverlessRoutesTable, plan string) (string, error) {
	pipe := m.Reader.TxPipeline()

	var getSessionRoute *redis.StringCmd
	if len(sessionID) > 0 {
		// Add the session to the rate limit table
		getSessionRoute = pipe.Get(ctx, getOrgSessionIDKey(orgID, serverlessRoutesTable, sessionID))
	}
	// checks user's active sessions against max allowed
	activeCount := pipe.ZCount(ctx, getOrgActiveServerlessCountKey(orgID, serverlessRoutesTable), fmt.Sprintf("%d", time.Now().Unix()), "inf")
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.Err(err).Msg("IrisCache: CheckServerlessSessionRateLimit: failed to get active sessions ZCount")
		return "", err
	}

	// this is returning the session route if it exists, as it has first priority
	getSessionRouteResult, rerr := getSessionRoute.Result()
	if rerr != nil && rerr != redis.Nil {
		log.Err(rerr).Msg("IrisCache: CheckServerlessSessionRateLimit: getSessionRouteResult failed to get session route")
		return "", rerr
	}
	// if it returns a valid route, then the session is already active, so we bypass the rate limit check
	if len(getSessionRouteResult) > 0 {
		return getSessionRouteResult, nil
	}

	if activeCount == nil || err == redis.Nil {
		return "", nil
	}
	// now checking if the user is rate limited
	activeCountResult, rerr := activeCount.Result()
	if rerr != nil && rerr != redis.Nil {
		log.Err(rerr).Msg("CheckServerlessSessionRateLimit: activeCountResult failed to get active session count")
		return "", rerr
	}
	if activeCountResult < 0 && rerr != nil {
		log.Err(rerr).Msg("CheckServerlessSessionRateLimit: error getting active session count")
		return "", rerr
	}

	// plan
	if int(activeCountResult) >= GetMonthlyPlanMaxAnvilServerlessSessions(plan) {
		err = fmt.Errorf("%d serverless %s sessions in use. max active sessions reached", activeCountResult, serverlessRoutesTable)
		log.Err(err).Msgf("GetNextServerlessRoute orgID: %d", orgID)
		return "", err
	}
	return "", nil
}

func (m *IrisCache) GetOrgActiveSessionsCount(ctx context.Context, orgID int, serverlessRoutesTable string) (int, error) {
	pipe := m.Reader.TxPipeline()

	// checks user's active sessions against max allowed
	activeCount := pipe.ZCount(ctx, getOrgActiveServerlessCountKey(orgID, serverlessRoutesTable), fmt.Sprintf("%d", time.Now().Unix()), "inf")
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return -1, err
	}

	activeCountResult, err := activeCount.Result()
	// redis == nil means there are no active sessions, so when that happens we skip and the implied zero return
	if err != nil && err != redis.Nil {
		log.Err(err)
		return -1, err
	}

	return int(activeCountResult), nil
}

/*
GetNextServerlessRoute returns the next available route from the serverless routing table in priority order
 1. Checks if the session is already active
 2. Checks if the session is rate limited
 3. Checks if there are any available routes
*/
func (m *IrisCache) GetNextServerlessRoute(ctx context.Context, orgID int, sessionID, serverlessRoutesTable, plan string) (string, bool, error) {
	path, err := m.CheckServerlessSessionRateLimit(ctx, orgID, sessionID, serverlessRoutesTable, plan)
	if err != nil {
		log.Err(err).Msg("GetNextServerlessRoute: CheckServerlessSessionRateLimit failed to check rate limit")
		return "", false, err
	}
	if len(path) > 0 {
		return path, false, nil
	}

	// tries to get a healthy route 3 times before failing
	for i := 0; i < 3; i++ {
		// done rate limit checks
		serviceTable := getGlobalServerlessTableKey(serverlessRoutesTable)
		pipe := m.Writer.TxPipeline()

		// Pop the first item from the list
		routeCmd := pipe.SPop(ctx, serviceTable)
		_, err = pipe.Exec(ctx)
		if err != nil {
			if err == redis.Nil {
				// error code, server unavailable or something tbd
				return "", false, fmt.Errorf("GetNextServerlessRoute: failed to find any available routes")
			}
			return "", false, err
		}
		path, err = routeCmd.Result()
		if err != nil {
			// error code, server unavailable or something tbd
			log.Err(err).Msg("GetNextServerlessRoute: failed to pop route")
			return "", false, err
		}
		if len(path) <= 0 {
			// error code, server unavailable or something tbd
			return "", false, fmt.Errorf("GetNextServerlessRoute: failed to find any available routes")
		}
		pipe = m.Writer.TxPipeline()
		redisSet := redis.Z{
			Score:  float64(time.Now().Add(ServerlessMaxRunTime).Unix()),
			Member: path,
		}
		// XX: Only update elements that already exist. Don't add new elements
		pipe.ZAddXX(ctx, getGlobalServerlessAvailabilityTableKey(serverlessRoutesTable), redisSet)

		// NX: Set key to hold string value if key does not exist.
		// this is used to k->v lookup for org-sessionID -> route
		pipe.SetNX(ctx, getOrgSessionIDKey(orgID, serverlessRoutesTable, sessionID), path, ServerlessSessionMaxRunTime)

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
			return "", false, err
		}
		resp, rerr := http.Get(path + "/health")
		if rerr == nil && resp.StatusCode < 400 && len(path) > 0 {
			return path, true, nil
		}
		if rerr != nil {
			path = ""
			rerr = nil
			log.Warn().Int("loopIndex", i).Msg("GetNextServerlessRoute: failed to get healthy route, retrying")
		}
		if resp != nil {
			log.Warn().Int("loopIndex", i).Int("statusCode", resp.StatusCode).Msg("GetNextServerlessRoute: failed to get healthy route, retrying")
		}

	}
	if len(path) <= 0 {
		return "", false, fmt.Errorf("GetNextServerlessRoute: failed to find any available routes")
	}
	return path, true, nil
}

func (m *IrisCache) GetServerlessSessionRoute(ctx context.Context, orgID int, serverlessRoutesTable, sessionID string) (string, error) {
	pipe := m.Reader.TxPipeline()
	// Gets and removes session-to-route cache key
	getSessionRoute := pipe.Get(ctx, getOrgSessionIDKey(orgID, serverlessRoutesTable, sessionID))

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msg("GetServerlessSessionRoute: failed to get session route")
		return "", err
	}
	// Get the value from the result of the Get command
	path, err := getSessionRoute.Result()
	if err != nil {
		log.Err(err).Msg("GetServerlessSessionRoute: failed to get session route")
		return "", err
	}
	return path, err
}

func (m *IrisCache) ReleaseServerlessRoute(ctx context.Context, orgID int, sessionID, serverlessRoutesTable string) (string, error) {
	// Get the value from the result of the Get command
	path, err := m.GetServerlessSessionRoute(ctx, orgID, serverlessRoutesTable, sessionID)
	if err != nil && err != redis.Nil {
		log.Err(err).Str("serverlessRoutesTable", serverlessRoutesTable).Str("sessionID", sessionID).Msg("ReleaseServerlessRoute: failed to get session route")
		return "", err
	}
	pipe := m.Writer.TxPipeline()
	// deletes session -> route cache key
	pipe.Del(ctx, getOrgSessionIDKey(orgID, serverlessRoutesTable, sessionID))

	// Resets the timing of the route with some margin
	if len(path) > 0 {
		// this is the user specific availability table
		sessionRateLimit := getOrgActiveServerlessCountKey(orgID, serverlessRoutesTable)
		pipe.ZRem(ctx, sessionRateLimit, path)

		// this is the global availability table
		// XX: Update elements that already exist. Don't add new elements
		serviceAvailabilityTable := getGlobalServerlessAvailabilityTableKey(serverlessRoutesTable)
		redisSet := redis.Z{
			Score:  float64(time.Now().Add(TimeMarginBuffer).Unix()),
			Member: path,
		}
		pipe.ZAddXX(ctx, serviceAvailabilityTable, redisSet)

		// this adds the route back to the set of available routes
		serverlessReadyRoutes := getGlobalServerlessTableKey(serverlessRoutesTable)
		pipe.SAdd(ctx, serverlessReadyRoutes, path)
	}

	// Removes expired sessions
	sessionRateLimit := getOrgActiveServerlessCountKey(orgID, serverlessRoutesTable)
	pipe.ZRemRangeByScore(ctx, sessionRateLimit, "0", fmt.Sprintf("%d", time.Now().Unix()))

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msg("ReleaseServerlessRoute: failed to release route")
		return "", err
	}

	return path, nil
}

func extractPodName(s string) (string, error) {
	// Parse the URL
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	// Split the host by "." to get the first part
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) > 0 {
		return parts[0], nil
	}

	return "", fmt.Errorf("cannot extract pod name from url: %s", s)
}

// Currently used for debug

func (m *IrisCache) GetServerlessRoutes(ctx context.Context, serverlessTableName string) ([]string, error) {
	pipe := m.Reader.TxPipeline()
	// Gets and removes session-to-route cache key

	serviceTable := getGlobalServerlessTableKey(serverlessTableName)
	getSessionRoutes := pipe.SMembers(ctx, serviceTable)
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msg("GetServerlessSessionRoute: failed to get session route")
		return nil, err
	}
	// Get the value from the result of the Get command
	path, err := getSessionRoutes.Result()
	if err != nil {
		log.Err(err).Msg("GetServerlessSessionRoute: failed to get session route")
		return nil, err
	}
	return path, err
}

func (m *IrisCache) DeleteAllServerlessTableArtifacts(ctx context.Context, serverlessTableName string) error {
	pipe := m.Writer.TxPipeline()
	// Gets and removes session-to-route cache key
	serviceTable := getGlobalServerlessTableKey(serverlessTableName)
	serviceCalendar := getGlobalServerlessAvailabilityTableKey(serverlessTableName)
	pipe.Del(ctx, serviceTable)
	pipe.Del(ctx, serviceCalendar)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msg("GetServerlessSessionRoute: failed to get session route")
		return err
	}

	return err
}

func (m *IrisCache) DeleteAllServerlessTableAndOrgArtifacts(ctx context.Context, orgID int, serverlessTableName, sessionID string) error {
	pipe := m.Writer.TxPipeline()
	// Gets and removes session-to-route cache key

	serviceTable := getGlobalServerlessTableKey(serverlessTableName)
	serviceCalendar := getGlobalServerlessAvailabilityTableKey(serverlessTableName)
	orgSession := getOrgSessionIDKey(orgID, serverlessTableName, sessionID)
	orgSessionCount := getOrgActiveServerlessCountKey(orgID, serverlessTableName)
	pipe.Del(ctx, serviceTable)
	pipe.Del(ctx, serviceCalendar)
	pipe.Del(ctx, orgSession)
	pipe.Del(ctx, orgSessionCount)

	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msg("GetServerlessSessionRoute: failed to get session route")
		return err
	}

	return err
}

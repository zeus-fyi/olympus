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
	ServerlessMaxRunTime = 10 * time.Minute
)

func (m *IrisCache) AddOrUpdateServerlessRoutingTable(ctx context.Context, serverlessRoutesTable string, routes []iris_models.RouteInfo) error {
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

func (m *IrisCache) GetNextServerlessRoute(ctx context.Context, serverlessRoutesTable string) (string, error) {
	//serviceTable := getGlobalServerlessTableKey(serverlessRoutesTable)
	//pipe := m.Writer.TxPipeline()

	return "", fmt.Errorf("GetNextServerlessRoute: failed to find any available routes")
}

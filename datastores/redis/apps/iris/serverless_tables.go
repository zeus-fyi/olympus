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
	serviceTable := getGlobalServerlessTableKey(serverlessRoutesTable)
	pipe := m.Writer.TxPipeline()
	for _, r := range routes {
		redisSet := redis.Z{
			Score:  float64(time.Now().Unix() + int64(ServerlessMaxRunTime.Seconds())),
			Member: r.RoutePath,
		}
		pipe.ZAddNX(ctx, serviceTable, redisSet)
	}
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return err
	}
	return nil
}

func (m *IrisCache) GetNextServerlessRoute(ctx context.Context, serverlessRoutesTable string) (string, error) {
	serviceTable := getGlobalServerlessTableKey(serverlessRoutesTable)
	pipe := m.Reader.TxPipeline()
	minElemCmd := pipe.ZRangeByScoreWithScores(ctx, serviceTable, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   fmt.Sprintf("%d", time.Now().Unix()),
		Count: 1, // Get only the first route
	})
	// Execute the transaction
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err)
		return "", err
	}
	member, err := minElemCmd.Result()
	if err != nil {
		log.Err(err).Msgf("GetAdaptiveEndpointByPriorityScoreAndInsertIfMissing")
		return "", err
	}
	if len(member) > 0 {
		path, eok := member[0].Member.(string)
		if !eok {
			return "", fmt.Errorf("GetNextServerlessRoute: failed to convert member to string")
		}
		if len(path) != 0 {
			return path, nil
		}
	}
	return "", fmt.Errorf("GetNextServerlessRoute: failed to find any available routes")
}

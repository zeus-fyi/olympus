package iris_redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (m *IrisCache) AddOrUpdateServerlessRoutingTable(ctx context.Context, serverlessRoutesTable string, routes []iris_models.RouteInfo) error {
	serviceTable := getGlobalServerlessTableKey(serverlessRoutesTable)
	pipe := m.Writer.TxPipeline()
	for _, r := range routes {
		redisSet := redis.Z{
			Score:  float64(time.Now().Unix()),
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

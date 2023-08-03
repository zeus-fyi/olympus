package iris_redis

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

var IrisRedisClient IrisCache

type IrisCache struct {
	Writer *redis.Client
	Reader *redis.Client
}

func NewIrisCache(ctx context.Context, w, r *redis.Client) IrisCache {
	log.Ctx(ctx).Info().Msg("IrisCache")
	log.Info().Interface("redis", r)
	return IrisCache{w, r}
}

func InitProductionRedisIrisCache(ctx context.Context) {
	writeRedisOpts := redis.Options{
		Addr: "redis-master.redis.svc.cluster.local:6379",
	}
	writer := redis.NewClient(&writeRedisOpts)
	readRedisOpts := redis.Options{
		Addr: "redis-replicas.redis.svc.cluster.local:6379",
	}
	reader := redis.NewClient(&readRedisOpts)
	IrisRedisClient = NewIrisCache(ctx, writer, reader)
}

func InitLocalTestProductionRedisIrisCache(ctx context.Context) {
	writeRedisOpts := redis.Options{
		Addr: "localhost:6379",
	}
	writer := redis.NewClient(&writeRedisOpts)
	readRedisOpts := redis.Options{
		Addr: "localhost:6380",
	}
	reader := redis.NewClient(&readRedisOpts)
	IrisRedisClient = NewIrisCache(ctx, writer, reader)
}

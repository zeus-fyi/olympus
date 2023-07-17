package iris_redis

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

type IrisCache struct {
	Writer *redis.Client
	Reader *redis.Client
}

func NewIrisCache(ctx context.Context, w, r *redis.Client) IrisCache {
	log.Ctx(ctx).Info().Msg("IrisCache")
	log.Info().Interface("redis", r)
	return IrisCache{w, r}
}

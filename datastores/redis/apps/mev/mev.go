package redis_mev

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

type MevCache struct {
	*redis.Client
}

func NewMevCache(ctx context.Context, r *redis.Client) MevCache {
	log.Ctx(ctx).Info().Msg("NewMevCache")
	log.Info().Interface("redis", r)
	return MevCache{r}
}

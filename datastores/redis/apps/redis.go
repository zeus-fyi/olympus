package redis_app

import (
	"context"

	"github.com/go-redis/redis/v9"
)

func InitRedis(ctx context.Context, opts redis.Options) *redis.Client {
	rdb := redis.NewClient(&opts)
	return rdb
}

package test_suites

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/configs"
	rdb "github.com/zeus-fyi/olympus/datastores/redis/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type RedisTestSuite struct {
	test_suites_base.TestSuite
	Redis *redis.Client
}

func (r *RedisTestSuite) SetupTest() {
	r.SetupRedisConn()
}

func (r *RedisTestSuite) SetupRedisConn() {
	r.Tc = configs.InitLocalTestConfigs()

	ctx := context.Background()
	redisOpts := redis.Options{}

	switch r.Tc.Env {
	case "local":
		redisOpts.Addr = r.Tc.LocalRedisConn
	case "staging":
		redisOpts.Addr = r.Tc.StagingRedisConn
	case "production":
	default:
		redisOpts.Addr = r.Tc.LocalRedisConn
	}

	r.Redis = rdb.InitRedis(ctx, redisOpts)
}

func (r *RedisTestSuite) SetupRedisConnStaging() {
	r.Tc = configs.InitLocalTestConfigs()
	ctx := context.Background()
	redisOpts := redis.Options{}
	redisOpts.Addr = r.Tc.StagingRedisConn
	r.Redis = rdb.InitRedis(ctx, redisOpts)
}

func (r *RedisTestSuite) SetupRedisConnProduction() {
	r.Tc = configs.InitProductionConfigs()
	ctx := context.Background()
	redisOpts := redis.Options{}
	redisOpts.Addr = r.Tc.LocalRedisConn
	r.Redis = rdb.InitRedis(ctx, redisOpts)
}

func (r *RedisTestSuite) CleanCache(ctx context.Context, keysToClear []string) {
	switch r.Tc.Env {
	case "local":
	case "staging":
		log.Info().Msg("not a local database, CleanupDb should only be used on a local database")
		return
	case "production":
		log.Info().Msg("not a local database, CleanupDb should only be used on a local database")
		return
	default:
	}

	for _, key := range keysToClear {
		r.Redis.Del(ctx, key)
	}
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}

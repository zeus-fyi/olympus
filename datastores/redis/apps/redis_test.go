package redis_app

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type RedisTestSuite struct {
	test_suites_base.TestSuite
}

func (r *RedisTestSuite) TestRedisConnection() {
	ctx := context.Background()
	redisOpts := redis.Options{
		Network: "",
		Addr:    "localhost:6379",
	}
	rdb := InitRedis(ctx, redisOpts)

	ttl := time.Second
	err := rdb.Set(ctx, "key", "value", ttl).Err()
	r.Require().Nil(err)
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}

package redis_mev

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	rdb "github.com/zeus-fyi/olympus/datastores/redis/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type BeaconIndexerCacheTestSuite struct {
	test_suites_base.TestSuite
	Redis *redis.Client
}

func (r *BeaconIndexerCacheTestSuite) TestTxCache() {
	ctx := context.Background()
	r.InitLocalConfigs()
	redisOpts := redis.Options{}
	redisOpts.Addr = r.Tc.LocalRedisConn
	r.Redis = rdb.InitRedis(ctx, redisOpts)
	apps.Pg.InitPG(ctx, r.Tc.LocalDbPgconn)

	fc := NewMevCache(ctx, r.Redis)
	r.Assert().NotEmpty(fc)

	ttl := time.Minute
	txHash := "tx"
	err := fc.AddTxHashCache(ctx, "tx", ttl)
	r.Require().Nil(err)

	doesKeyExist, err := fc.DoesTxExist(ctx, txHash)
	r.Require().Nil(err)
	r.Assert().True(doesKeyExist)

	err = fc.DeleteTx(ctx, txHash)
	r.Require().Nil(err)
}

func TestBeaconIndexerCacheTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconIndexerCacheTestSuite))
}

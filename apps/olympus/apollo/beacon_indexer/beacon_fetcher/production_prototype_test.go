package beacon_fetcher

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	rdb "github.com/zeus-fyi/olympus/datastores/redis/apps"
	"github.com/zeus-fyi/olympus/datastores/redis/apps/beacon_indexer"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	apollo_buckets "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/buckets"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type ProductionBeaconFetcherTestSuite struct {
	test_suites_base.TestSuite
	Redis *redis.Client
}

func (f *ProductionBeaconFetcherTestSuite) SetupTest() {
	ctx := context.Background()
	f.InitProductionConfig()

	redisOpts := redis.Options{}
	redisOpts.Addr = f.Tc.ProdRedisConn
	f.Redis = rdb.InitRedis(ctx, redisOpts)
	apps.Pg.InitPG(ctx, f.Tc.ProdLocalApolloDbPgconn)
	tc := configs.InitLocalTestConfigs()
	authKeysCfg := tc.ProdLocalAuthKeysCfg
	apollo_buckets.ApolloS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)
}

func (f *ProductionBeaconFetcherTestSuite) TestInsertNewEpochCheckpoint() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)
	InsertCheckpointsTimeout = time.Second
	InsertNewEpochCheckpoint()
}
func (f *ProductionBeaconFetcherTestSuite) TestFetchAndCacheAnyAfterCheckpoint() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)
	err := fetchAnyValidatorBalancesAfterCheckpoint(ctx, 5*time.Minute)
	f.Require().Nil(err)
}
func (f *ProductionBeaconFetcherTestSuite) TestFetchNewOrMissingValidators() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)

	checkpointEpoch = 163999
	err := Fetcher.BeaconFindNewAndMissingValidatorIndexes(ctx, 10)
	f.Require().Nil(err)
}

func (f *ProductionBeaconFetcherTestSuite) TestInsertForwardFetchCheckpointUsingCache() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)

	checkpointEpoch = 163999
	err := Fetcher.InsertForwardFetchCheckpointUsingCache(ctx)
	f.Require().Nil(err)
}

func TestProductionBeaconFetcherTestSuite(t *testing.T) {
	suite.Run(t, new(ProductionBeaconFetcherTestSuite))
}

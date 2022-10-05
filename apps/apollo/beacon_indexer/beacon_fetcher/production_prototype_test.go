package beacon_fetcher

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/redis_app/beacon_indexer"
	"github.com/zeus-fyi/olympus/pkg/utils/env"
)

type ProductionBeaconFetcherTestSuite struct {
	env.ProductionPrototypeTest
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

	checkpointEpoch = 138850
	err := Fetcher.BeaconFindNewAndMissingValidatorIndexes(ctx, 500)
	f.Require().Nil(err)
}

func (f *ProductionBeaconFetcherTestSuite) TestInsertForwardFetchCheckpointUsingCache() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)

	checkpointEpoch = 138850
	err := Fetcher.InsertForwardFetchCheckpointUsingCache(ctx)
	f.Require().Nil(err)
}

func TestProductionBeaconFetcherTestSuite(t *testing.T) {
	suite.Run(t, new(ProductionBeaconFetcherTestSuite))
}

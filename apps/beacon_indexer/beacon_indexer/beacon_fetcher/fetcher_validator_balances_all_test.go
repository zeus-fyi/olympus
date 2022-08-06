package beacon_fetcher

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
	"github.com/zeus-fyi/olympus/pkg/datastores/redis_app/beacon_indexer"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type BeaconFetcherTestSuite struct {
	test_suites.DatastoresTestSuite
}

func (f *BeaconFetcherTestSuite) TestCheckpointCache() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)

	chkPoint := beacon_models.ValidatorsEpochCheckpoint{}
	err := chkPoint.GetFirstEpochCheckpointWithBalancesRemaining(ctx)
	f.Require().Nil(err)
	doesExist, err := Fetcher.Cache.DoesCheckpointExist(ctx, chkPoint.Epoch)
	f.Require().Nil(err)
	f.Assert().False(doesExist)

	err = fetchAllValidatorBalances(ctx, 5*time.Minute)

	doesExist, err = Fetcher.Cache.DoesCheckpointExist(ctx, chkPoint.Epoch)
	f.Require().Nil(err)
	f.Assert().True(doesExist)

}

func (f *BeaconFetcherTestSuite) TestFetchAndCacheAnyAfterCheckpoint() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)
	err := fetchAnyValidatorBalancesAfterCheckpoint(ctx, 5*time.Minute)
	f.Require().Nil(err)
}

func TestBeaconFetcherTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconFetcherTestSuite))
}

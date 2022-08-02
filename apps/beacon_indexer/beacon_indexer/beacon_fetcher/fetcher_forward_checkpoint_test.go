package beacon_fetcher

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/datastores/redis_app/beacon_indexer"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type BeaconForwardCheckpointFetcherTestSuite struct {
	test_suites.DatastoresTestSuite
}

func (f *BeaconFetcherTestSuite) TestForwardFetchCheckpoint() {
	ctx := context.Background()
	fetcher.NodeEndpoint = f.Tc.LocalBeaconConn

	epoch := int64(134000)
	vbForwardCheckpoint, err := fetcher.FetchForwardCheckpointValidatorBalances(ctx, epoch)
	f.Require().Nil(err)
	f.Assert().NotEmpty(vbForwardCheckpoint.ValidatorBalances)

	prevBalanceEpochMap := map[int64]int64{}
	prevEpoch, err := fetcher.FetchAllValidatorBalances(ctx, epoch-1)

	f.Require().Nil(err)
	f.Assert().NotEmpty(prevEpoch.ValidatorBalances)

	for _, balance := range prevEpoch.ValidatorBalances {
		prevBalanceEpochMap[balance.Index] = balance.TotalBalanceGwei
	}

	currentEpoch, err := fetcher.FetchAllValidatorBalances(ctx, epoch)
	f.Require().Nil(err)
	f.Assert().NotEmpty(currentEpoch.ValidatorBalances)
	for _, balance := range currentEpoch.ValidatorBalances {
		if _, ok := prevBalanceEpochMap[balance.Index]; ok {
			prevBalanceEpochMap[balance.Index] = balance.TotalBalanceGwei - prevBalanceEpochMap[balance.Index]
		} else {
			prevBalanceEpochMap[balance.Index] = 0
		}
	}

	for _, val := range vbForwardCheckpoint.ValidatorBalances {
		f.Assert().Equal(prevBalanceEpochMap[val.Index], val.CurrentEpochYieldGwei)
	}

	err = vbForwardCheckpoint.InsertValidatorBalances(ctx)
	f.Assert().Nil(err)
}

func (f *BeaconFetcherTestSuite) TestForwardCheckpointBalanceUpdate() {
	ctx := context.Background()
	fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)
	err := fetchAllValidatorBalancesAfterCheckpoint(ctx, 10*time.Minute)
	f.Require().Nil(err)
}

func (f *BeaconFetcherTestSuite) TestCache() {
	ctx := context.Background()
	fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)
	epoch := 5

	key, err := fetcher.Cache.SetCheckpointCache(ctx, epoch, time.Minute)
	f.Require().Nil(err)
	f.Assert().NotEmpty(key)
	doesExist, err := fetcher.Cache.DoesCheckpointExist(ctx, epoch)
	f.Require().Nil(err)
	f.Assert().True(doesExist)
}
func TestBeaconForwardCheckpointFetcherTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconForwardCheckpointFetcherTestSuite))
}

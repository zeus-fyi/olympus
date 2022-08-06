package beacon_fetcher

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/datastores/redis_app/beacon_indexer"
	"github.com/zeus-fyi/olympus/pkg/utils/env"
)

type StagingBeaconFetcherTestSuite struct {
	env.StagingPrototypeTest
}

func (f *StagingBeaconFetcherTestSuite) TestSimulateAppOnStaging() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn

	InitFetcherService(ctx, f.Tc.LocalBeaconConn, f.Redis)
	checkpointEpoch = 1

	log.Info().Msg("will not return here")
	for {
		time.Sleep(10 * time.Minute)
	}
}

func (f *StagingBeaconFetcherTestSuite) TestFetchAndCacheAnyAfterCheckpoint() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)
	err := fetchAnyValidatorBalancesAfterCheckpoint(ctx, 5*time.Minute)
	f.Require().Nil(err)
}

func (f *StagingBeaconFetcherTestSuite) TestInsertForwardFetchCheckpoint() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)

	epoch := int64(10)
	vbForwardCheckpoint, err := Fetcher.FetchForwardCheckpointValidatorBalances(ctx, epoch)
	f.Require().Nil(err)
	f.Assert().NotEmpty(vbForwardCheckpoint.ValidatorBalances)

	prevBalanceEpochMap := map[int64]int64{}
	prevEpoch, err := Fetcher.FetchAllValidatorBalances(ctx, epoch-1)

	f.Require().Nil(err)
	f.Assert().NotEmpty(prevEpoch.ValidatorBalances)

	for _, balance := range prevEpoch.ValidatorBalances {
		prevBalanceEpochMap[balance.Index] = balance.TotalBalanceGwei
	}

	currentEpoch, err := Fetcher.FetchAllValidatorBalances(ctx, epoch)
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

func (f *StagingBeaconFetcherTestSuite) TestCache() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)

	epoch := int64(10)
	vbForwardCheckpoint, err := Fetcher.FetchForwardCheckpointValidatorBalances(ctx, epoch)
	f.Require().Nil(err)
	_, err = Fetcher.Cache.SetBalanceCache(ctx, int(epoch), vbForwardCheckpoint, time.Hour*24*7)
	f.Require().Nil(err)

	expVbeCache, err := Fetcher.Cache.GetBalanceCache(ctx, int(epoch))
	f.Require().Nil(err)
	f.Assert().NotEmpty(expVbeCache)
}

func (f *StagingBeaconFetcherTestSuite) TestInsertForwardFetchCheckpointUsingCache() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)

	checkpointEpoch = 1
	err := Fetcher.InsertForwardFetchCheckpointUsingCache(ctx)
	f.Require().Nil(err)
}

func TestStagingBeaconFetcherTestSuite(t *testing.T) {
	suite.Run(t, new(StagingBeaconFetcherTestSuite))
}

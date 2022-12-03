package beacon_fetcher

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/redis/apps/beacon_indexer"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type BeaconForwardCheckpointFetcherTestSuite struct {
	test_suites.DatastoresTestSuite
}

func (f *BeaconFetcherTestSuite) TestCache() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn
	Fetcher.Cache = beacon_indexer.NewFetcherCache(ctx, f.Redis)
	epoch := 5

	key, err := Fetcher.Cache.SetCheckpointCache(ctx, epoch, time.Minute)
	f.Require().Nil(err)
	f.Assert().NotEmpty(key)
	doesExist, err := Fetcher.Cache.DoesCheckpointExist(ctx, epoch)
	f.Require().Nil(err)
	f.Assert().True(doesExist)
}

func TestBeaconForwardCheckpointFetcherTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconForwardCheckpointFetcherTestSuite))
}

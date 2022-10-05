package beacon_indexer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type BeaconIndexerCacheTestSuite struct {
	test_suites.DatastoresTestSuite
}

func (r *BeaconIndexerCacheTestSuite) TestCheckpointCache() {
	ctx := context.Background()
	fc := NewFetcherCache(ctx, r.Redis)
	r.Assert().NotEmpty(fc)

	epoch := 1
	ttl := time.Minute
	key, err := fc.SetCheckpointCache(ctx, epoch, ttl)
	r.Require().Nil(err)
	chkPoint := fc.Get(ctx, key)
	val, err := chkPoint.Int()
	r.Require().Nil(err)
	r.Assert().Equal(epoch, val)

	doesKeyExist, err := fc.DoesCheckpointExist(ctx, epoch)
	r.Require().Nil(err)
	r.Assert().True(doesKeyExist)

	doesKeyExist, err = fc.DoesCheckpointExist(ctx, epoch+1)
	r.Require().Nil(err)
	r.Assert().False(doesKeyExist)

	err = fc.DeleteCheckpoint(ctx, epoch)
	r.Require().Nil(err)
}

func TestBeaconIndexerCacheTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconIndexerCacheTestSuite))
}

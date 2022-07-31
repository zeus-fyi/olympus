package beacon_indexer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type BeaconIndexerCacheTestSuite struct {
	test_suites.DatastoresTestSuite
}

func (r *BeaconIndexerCacheTestSuite) TestCheckpointCache() {
	ctx := context.Background()
	fc := FetcherCache{r.Redis}
	r.Assert().NotEmpty(fc)

	fc.CheckpointCache(ctx)

	chkPoint := fc.Get(ctx, "k")
	r.Assert().NotEmpty(chkPoint)
}

func TestBeaconIndexerCacheTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconIndexerCacheTestSuite))
}

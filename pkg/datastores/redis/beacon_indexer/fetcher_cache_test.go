package beacon_indexer

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type BeaconIndexerCacheTestSuite struct {
	test_suites.DatastoresTestSuite
}

func (r *BeaconIndexerCacheTestSuite) TestCheckpointCache() {

}

func TestBeaconIndexerCacheTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconIndexerCacheTestSuite))
}

package beacon_fetcher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type BeaconStatusUpdateFetcherTestSuite struct {
	test_suites.DatastoresTestSuite
}

func (f *BeaconFetcherTestSuite) TestFetcherUpdateAll() {
	ctx := context.Background()
	fetcher.NodeEndpoint = f.Tc.LocalBeaconConn

	err := fetcher.BeaconUpdateAllValidatorStates(ctx)
	f.Require().Nil(err)

	f.Assert().NotEmpty(fetcher.Validators)
}

func TestBeaconStatusUpdateFetcherTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconStatusUpdateFetcherTestSuite))
}

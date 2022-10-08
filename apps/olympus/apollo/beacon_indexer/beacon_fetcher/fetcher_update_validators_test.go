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
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn

	err := Fetcher.BeaconUpdateAllValidatorStates(ctx)
	f.Require().Nil(err)

	f.Assert().NotEmpty(Fetcher.Validators)
}

func (f *BeaconFetcherTestSuite) TestFetcherUpdateBatch() {
	ctx := context.Background()
	Fetcher.NodeEndpoint = f.Tc.LocalBeaconConn

	err := Fetcher.BeaconUpdateValidatorStates(ctx, 100)
	f.Require().Nil(err)

	f.Assert().NotEmpty(Fetcher.Validators)
}

func TestBeaconStatusUpdateFetcherTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconStatusUpdateFetcherTestSuite))
}

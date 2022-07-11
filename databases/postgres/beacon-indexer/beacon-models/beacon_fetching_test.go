package beacon_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BeaconFetchingTestSuite struct {
	BeaconBaseTestSuite
}

func (f *BeaconFetchingTestSuite) TestSelectValidatorsFromBeaconAPI() {
	v, err := SelectValidatorsToQueryBeaconForBalanceUpdates(context.Background(), 100)
	f.Require().Nil(err)
	f.Assert().NotEmpty(v)
}

func (f *BeaconFetchingTestSuite) TestSelectValidatorIndexesInStrArrayForQueryURL() {
	v, err := SelectValidatorsQueryOngoingStates(context.Background(), 100)
	f.Require().Nil(err)
	f.Assert().NotEmpty(v)
}

func TestBeaconFetchingTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconFetchingTestSuite))
}

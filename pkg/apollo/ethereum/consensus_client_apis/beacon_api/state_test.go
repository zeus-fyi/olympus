package beacon_api

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BeaconStateAPITestSuite struct {
	BeaconAPITestSuite
}

func (s *BeaconStateAPITestSuite) TestGetValidatorsByState() {
	s.SkipTest(disableHighDataAPITests)
	state := "finalized"
	status := "active_ongoing"
	var vs ValidatorsStateBeacon
	vs, err := vs.FetchAllStateAndDecode(ctx, s.Tc.LocalBeaconConn, state, status)
	s.Require().Nil(err)
	s.Assert().NotEmpty(vs)

}

func TestBeaconStateAPITestSuite(t *testing.T) {
	suite.Run(t, new(BeaconStateAPITestSuite))
}

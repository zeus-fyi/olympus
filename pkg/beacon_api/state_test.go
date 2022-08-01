package beacon_api

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type BeaconStateAPITestSuite struct {
	base.TestSuite
}

func (s *BeaconStateAPITestSuite) TestGetValidatorsByState() {
	s.SkipTest(disableHighDataAPITests)
	state := "finalized"
	status := "active_ongoing"
	var vs ValidatorsStateBeacon
	err := vs.FetchAllStateAndDecode(ctx, s.Tc.LocalBeaconConn, state, status)
	s.Require().Nil(err)
	s.Assert().NotEmpty(vs)

}

func TestBeaconStateAPITestSuite(t *testing.T) {
	suite.Run(t, new(BeaconStateAPITestSuite))
}

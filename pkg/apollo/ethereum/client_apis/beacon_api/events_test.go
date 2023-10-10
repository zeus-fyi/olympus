package beacon_api

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BeaconEventsAPITestSuite struct {
	BeaconAPITestSuite
}

func (s *BeaconEventsAPITestSuite) TestGetValidatorsByState() {
	topic := "attestation"
	err := beaconSubscriptionSSE(ctx, s.Tc.QuikNodeLiveNode, topic)
	s.Require().Nil(err)
}

func TestBeaconEventsAPITestSuite(t *testing.T) {
	suite.Run(t, new(BeaconEventsAPITestSuite))
}

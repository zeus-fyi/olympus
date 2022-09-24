package beacon_api

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

const nodeURL = "https://sleek-hidden-frost.ethereum-goerli.discover.quiknode.pro/c19f411cb9561127df0126e4c57d959122794dcd"

type BeaconEventsAPITestSuite struct {
	BeaconAPITestSuite
}

func (s *BeaconEventsAPITestSuite) TestGetValidatorsByState() {
	topic := "attestation"
	err := beaconSubscriptionSSE(ctx, nodeURL, topic)
	s.Require().Nil(err)
}

func TestBeaconEventsAPITestSuite(t *testing.T) {
	suite.Run(t, new(BeaconEventsAPITestSuite))
}

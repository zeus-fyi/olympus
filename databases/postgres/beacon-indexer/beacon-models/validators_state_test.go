package beacon_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidatorsStateTestSuite struct {
	BeaconBaseTestSuite
}

func (s *ValidatorsStateTestSuite) TestInsertValidatorsFromBeaconAPI() {
	var vs Validators
	err := vs.InsertValidatorsFromBeaconAPI(context.Background())
	s.Require().Nil(err)
}

func TestValidatorsStateTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorsStateTestSuite))
}

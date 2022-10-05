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

func (s *ValidatorsStateTestSuite) TestUpdateValidatorsFromBeaconAPI() {
	var vs Validators
	ctx := context.Background()
	vs = createAndInsertFakeValidatorsByStatus(2, "active_ongoing_genesis")
	for i, _ := range vs.Validators {
		vs.Validators[i].WithdrawableEpoch = 122222222
	}
	rows, err := vs.UpdateValidatorsFromBeaconAPI(ctx)
	s.Require().Nil(err)
	s.Assert().Equal(int64(2), rows)
}

func TestValidatorsStateTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorsStateTestSuite))
}

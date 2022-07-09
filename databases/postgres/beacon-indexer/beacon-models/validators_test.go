package beacon_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidatorsTestSuite struct {
	BeaconBaseTestSuite
}

func (s *ValidatorsTestSuite) TestInsertAndSelectValidators() {
	vs := createFakeValidators(2)
	err := vs.InsertValidators(context.Background())
	s.Require().Nil(err)

	selectedVs, err := vs.SelectValidators(context.Background())
	s.Require().Nil(err)
	s.Assert().Len(selectedVs.Validators, 2)
}

func TestValidatorsTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorsTestSuite))
}

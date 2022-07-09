package beacon_models

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidatorBalancesTestSuite struct {
	BeaconBaseTestSuite
}

func (s *ValidatorBalancesTestSuite) TestInsertValidatorBalances() {

}

func TestValidatorBalancesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorBalancesTestSuite))
}

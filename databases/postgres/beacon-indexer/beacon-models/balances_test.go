package beacon_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidatorBalancesTestSuite struct {
	BeaconBaseTestSuite
}

func (s *ValidatorBalancesTestSuite) TestInsertAndSelectValidatorBalances() {
	vs := createFakeValidators(3)
	ctx := context.Background()
	err := vs.InsertValidators(ctx)
	s.Require().Nil(err)

	expVBs := make(map[int64]ValidatorBalanceEpoch, 3)

	var vbs ValidatorBalancesEpoch
	for i, v := range vs.Validators {
		var vb ValidatorBalanceEpoch
		vb.Epoch = int64(i)
		vb.Validator = v
		vb.TotalBalanceGwei = int64(32000000001 - i)
		vbs.ValidatorBalance = append(vbs.ValidatorBalance, vb)
		expVBs[vb.Index] = vb
	}

	err = vbs.InsertValidatorBalances(ctx)
	s.Require().Nil(err)

	selectedVBs, err := vbs.SelectValidatorBalances(ctx)
	s.Require().Nil(err)

	for _, testVB := range selectedVBs.ValidatorBalance {
		expVB := expVBs[testVB.Index]
		s.Require().NotEmpty(expVB)
		s.Assert().Equal(expVB.Index, testVB.Index)
		s.Assert().Equal(expVB.TotalBalanceGwei, testVB.TotalBalanceGwei, "validator_index %s total balance expected to be %s, but was %s", expVB.Index, expVB.TotalBalanceGwei, testVB.TotalBalanceGwei)
		s.Assert().Equal(expVB.TotalBalanceGwei-32000000000, testVB.YieldToDateGwei, "validator_index %s yield_to_date_gwei expected to be %s, but was %s", expVB.Index, expVB.YieldToDateGwei, testVB.TotalBalanceGwei)
		//s.Assert().Equal(expVB.CurrentEpochYieldGwei, testVB.CurrentEpochYieldGwei)
	}
}

func TestValidatorBalancesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorBalancesTestSuite))
}

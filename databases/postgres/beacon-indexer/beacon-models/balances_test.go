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
	}
}

func (s *ValidatorBalancesTestSuite) TestInsertValidatorBalancesForNextEpoch() {
	vs := createFakeValidators(2)
	ctx := context.Background()
	err := vs.InsertValidators(ctx)
	s.Require().Nil(err)

	expVBs := make(map[int64]ValidatorBalanceEpoch, 2)

	var vbsInitial ValidatorBalancesEpoch
	for _, v := range vs.Validators {
		var vb ValidatorBalanceEpoch
		vb.Epoch = int64(0)
		vb.Validator = v
		vb.TotalBalanceGwei = int64(32000000000)
		vbsInitial.ValidatorBalance = append(vbsInitial.ValidatorBalance, vb)
		expVBs[vb.Index] = vb
	}

	err = vbsInitial.InsertValidatorBalances(ctx)
	s.Require().Nil(err)

	selectedVBs, err := vbsInitial.SelectValidatorBalances(ctx)
	s.Require().Nil(err)

	for _, testVB := range selectedVBs.ValidatorBalance {
		expVB := expVBs[testVB.Index]
		s.Require().NotEmpty(expVB)
		s.Assert().Equal(expVB.Index, testVB.Index)
		s.Assert().Equal(expVB.TotalBalanceGwei, testVB.TotalBalanceGwei, "validator_index %s total balance expected to be %s, but was %s", expVB.Index, expVB.TotalBalanceGwei, testVB.TotalBalanceGwei)
		s.Assert().Equal(expVB.TotalBalanceGwei-32000000000, testVB.YieldToDateGwei, "validator_index %s yield_to_date_gwei expected to be %s, but was %s", expVB.Index, expVB.YieldToDateGwei, testVB.TotalBalanceGwei)
		s.Assert().Equal(int64(0), testVB.CurrentEpochYieldGwei, "validator_index %s yield_to_date_gwei expected to be %s, but was %s", expVB.Index, 0, testVB.TotalBalanceGwei)
	}

	var vbsNextEpoch ValidatorBalancesEpoch
	expVBs = make(map[int64]ValidatorBalanceEpoch, 2)
	for _, v := range vs.Validators {
		var vb ValidatorBalanceEpoch
		vb.Epoch = int64(1)
		vb.Validator = v
		vb.TotalBalanceGwei = int64(32000000010)
		vbsNextEpoch.ValidatorBalance = append(vbsNextEpoch.ValidatorBalance, vb)
		expVBs[vb.Index] = vb
	}

	err = vbsNextEpoch.InsertValidatorBalancesForNextEpoch(ctx)
	s.Require().Nil(err)

	selectedVBs, err = vbsNextEpoch.SelectValidatorBalances(ctx)
	s.Require().Nil(err)

	for _, testVB := range selectedVBs.ValidatorBalance {
		if testVB.Epoch > 0 {
			expVB := expVBs[testVB.Index]
			s.Require().NotEmpty(expVB)
			s.Assert().Equal(expVB.Index, testVB.Index)
			s.Assert().Equal(expVB.Epoch, testVB.Epoch)
			s.Assert().Equal(int64(1), testVB.Epoch)
			expBalance := int64(32000000010)
			expYieldDiff := int64(32000000010 - 32000000000)
			s.Assert().Equal(expBalance, testVB.TotalBalanceGwei, "validator_index %s total balance expected to be %s, but was %s", expVB.Index, expVB.TotalBalanceGwei, testVB.TotalBalanceGwei)
			s.Assert().Equal(expYieldDiff, testVB.YieldToDateGwei, "validator_index %s yield_to_date_gwei expected to be %s, but was %s", expVB.Index, expVB.YieldToDateGwei, testVB.TotalBalanceGwei)
			s.Assert().Equal(expYieldDiff, testVB.CurrentEpochYieldGwei, "validator_index %s yield_to_date_gwei expected to be %s, but was %s", expVB.Index, expYieldDiff, testVB.TotalBalanceGwei)
		}
	}
}

func TestValidatorBalancesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorBalancesTestSuite))
}

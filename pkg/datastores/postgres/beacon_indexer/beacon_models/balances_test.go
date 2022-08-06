package beacon_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidatorBalancesTestSuite struct {
	BeaconBaseTestSuite
}

func (s *ValidatorBalancesTestSuite) TestInsertValidatorBalancesNoChecks() {
	ctx := context.Background()

	vs := createAndInsertFakeValidatorsByStatus(2, "active_ongoing_genesis")

	vbe := ValidatorBalancesEpoch{}
	vbe.ValidatorBalances = make([]ValidatorBalanceEpoch, len(vs.Validators))
	for i, v := range vs.Validators {
		var vbAtEpoch ValidatorBalanceEpoch
		vbAtEpoch.Validator = v
		vbAtEpoch.Epoch = 134000
		vbAtEpoch.TotalBalanceGwei = eth32Gwei + 342355
		vbAtEpoch.CurrentEpochYieldGwei = 0
		vbe.ValidatorBalances[i] = vbAtEpoch
	}

	err := vbe.InsertValidatorBalances(ctx)
	s.Require().Nil(err)
}

func (s *ValidatorBalancesTestSuite) TestInsertValidatorBalancesForNextEpoch() {
	ctx := context.Background()
	vs := seedValidators(2)

	epoch0Balance := int64(32000000000)
	s.insertEpochVBFromGivenValidatorsAndValidate(ctx, vs, int64(0), epoch0Balance, epoch0Balance)

	epoch1Balance := int64(32000000010)
	s.insertEpochVBFromGivenValidatorsAndValidate(ctx, vs, int64(1), epoch1Balance, epoch0Balance)

	epoch2Balance := int64(32000000008)
	s.insertEpochVBFromGivenValidatorsAndValidate(ctx, vs, int64(2), epoch2Balance, epoch1Balance)
}

func (s *ValidatorBalancesTestSuite) TestInsertValidatorBalancesForAnyEpoch() {
	//ctx := context.Background()
	//vs := seedValidators(2)
	//
	//epoch0Balance := int64(32000000000)
	//vbsEpoch, _ := seedValidatorEpochBalancesStructFromValidators(vs, 100, epoch0Balance)
	//err := vbsEpoch.InsertValidatorBalancesForAnyEpoch(ctx)
	//s.Require().Nil(err)
}

func (s *ValidatorBalancesTestSuite) TestInsertValidatorBalancesForNextEpochEdgeCases() {
	ctx := context.Background()
	vs := seedValidators(1)

	// test when validator not in validator table
	for _, v := range vs.Validators {
		v.Index = int64(0)
	}

	epoch0Balance := int64(32000000000)
	vbsEpoch, _ := seedAndInsertNewValidatorBalances(ctx, vs, int64(0), epoch0Balance)
	selectedVbsEpoch0, _ := vbsEpoch.SelectValidatorBalances(ctx)
	s.Assert().Empty(selectedVbsEpoch0)

	// test when trying to add >1 epoch ahead not in validator table
	vs = seedValidators(2)
	vbsEpoch, _ = seedAndInsertNewValidatorBalances(ctx, vs, int64(3), epoch0Balance)
	selectedVbsEpoch3, _ := vbsEpoch.SelectValidatorBalances(ctx)
	s.Assert().Empty(selectedVbsEpoch3)
}

func (s *ValidatorBalancesTestSuite) TestInsertAndSelectValidatorBalances() {
	vs := seedValidators(3)
	ctx := context.Background()

	vbsEpoch0, exp0VBs := seedValidatorEpochBalancesStructFromValidators(vs, int64(0), int64(32000000000))
	err := vbsEpoch0.InsertValidatorBalances(ctx)
	s.Require().Nil(err)

	selectedVBs, err := vbsEpoch0.SelectValidatorBalances(ctx)
	s.Require().Nil(err)

	for _, testVB := range selectedVBs.ValidatorBalances {
		expVB := exp0VBs[testVB.Index]
		s.Require().NotEmpty(expVB)
		s.Assert().Equal(expVB.Index, testVB.Index)
		s.Assert().Equal(expVB.TotalBalanceGwei, testVB.TotalBalanceGwei, "validator_index %s total balance expected to be %s, but was %s", expVB.Index, expVB.TotalBalanceGwei, testVB.TotalBalanceGwei)
		s.Assert().Equal(expVB.TotalBalanceGwei-32000000000, testVB.YieldToDateGwei, "validator_index %s yield_to_date_gwei expected to be %s, but was %s", expVB.Index, expVB.YieldToDateGwei, testVB.TotalBalanceGwei)
	}
}

func (s *ValidatorBalancesTestSuite) insertEpochVBFromGivenValidatorsAndValidate(ctx context.Context, vs Validators, epoch, epochTotalBalance, prevEpochTotalBalance int64) {
	vbsEpoch, expVbsEpoch := seedAndInsertNewValidatorBalances(ctx, vs, epoch, epochTotalBalance)
	selectedVbsEpoch0, err := vbsEpoch.SelectValidatorBalances(ctx)
	s.Require().Nil(err)
	s.Require().NotNil(selectedVbsEpoch0)
	s.compareExpectedToActualValBalanceEntries(*selectedVbsEpoch0, expVbsEpoch, epochTotalBalance, prevEpochTotalBalance)
}

func (s *ValidatorBalancesTestSuite) compareExpectedToActualValBalanceEntries(selectedVBs ValidatorBalancesEpoch, expVBs map[int64]ValidatorBalanceEpoch, startingTotalBalance, prevTotalBalanceGwei int64) {
	for _, testVB := range selectedVBs.ValidatorBalances {
		expVB := expVBs[testVB.Index]
		s.Require().NotEmpty(expVB)
		s.Assert().Equal(expVB.Index, testVB.Index)
		// Total Balance
		s.Assert().Equal(expVB.TotalBalanceGwei, testVB.TotalBalanceGwei, "validator_index %s total balance expected to be %s, but was %s", expVB.Index, expVB.TotalBalanceGwei, testVB.TotalBalanceGwei)
		// Yield To Date Gwei
		yieldToDateGwei := expVB.TotalBalanceGwei - startingTotalBalance
		s.Assert().Equal(yieldToDateGwei, testVB.YieldToDateGwei, "validator_index %s yield_to_date_gwei expected to be %s, but was %s", expVB.Index, yieldToDateGwei, testVB.YieldToDateGwei)
		// Epoch Yield Gwei
		epochYieldGwei := expVB.TotalBalanceGwei - prevTotalBalanceGwei
		s.Assert().Equal(epochYieldGwei, testVB.CurrentEpochYieldGwei, "validator_index %s epoch_yield_gwei expected to be %s, but was %s", expVB.Index, epochYieldGwei, testVB.CurrentEpochYieldGwei)
	}
}

func seedAndInsertNewValidatorBalances(ctx context.Context, vs Validators, epoch, epochTotalBalance int64) (ValidatorBalancesEpoch, map[int64]ValidatorBalanceEpoch) {
	vbsEpoch, expVbsMapEpoch := seedValidatorEpochBalancesStructFromValidators(vs, epoch, epochTotalBalance)
	err := vbsEpoch.InsertValidatorBalancesForNextEpoch(ctx)
	if err != nil {
		panic(err)
	}

	return vbsEpoch, expVbsMapEpoch
}

func seedValidators(numValidators int) Validators {
	vs := createFakeValidators(numValidators)
	ctx := context.Background()
	err := vs.InsertValidatorsOnlyIndexPubkey(ctx)
	if err != nil {
		panic(err)
	}
	return vs
}

func seedValidatorEpochBalancesStructFromValidators(vs Validators, epoch, newGweiBalance int64) (ValidatorBalancesEpoch, map[int64]ValidatorBalanceEpoch) {
	var valExpBalancesAtEpoch ValidatorBalancesEpoch
	expVbsMap := make(map[int64]ValidatorBalanceEpoch, len(vs.Validators)-1)
	for _, v := range vs.Validators {
		var vb ValidatorBalanceEpoch
		vb.Epoch = epoch
		vb.Validator = v
		vb.TotalBalanceGwei = newGweiBalance
		valExpBalancesAtEpoch.ValidatorBalances = append(valExpBalancesAtEpoch.ValidatorBalances, vb)
		expVbsMap[vb.Index] = vb
	}
	return valExpBalancesAtEpoch, expVbsMap
}

func TestValidatorBalancesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorBalancesTestSuite))
}

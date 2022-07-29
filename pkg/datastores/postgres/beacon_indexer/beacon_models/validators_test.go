package beacon_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidatorsTestSuite struct {
	BeaconBaseTestSuite
}

func (s *ValidatorsTestSuite) TestInsertValidatorsFromBeaconAPI() {
	var vs Validators
	err := vs.InsertValidatorsFromBeaconAPI(context.Background())
	s.Require().Nil(err)
}

func (s *ValidatorsTestSuite) TestInsertAndSelectValidatorsIndexPubkey() {
	vs := createFakeValidators(2)
	err := vs.InsertValidatorsOnlyIndexPubkey(context.Background())
	s.Require().Nil(err)

	selectedVs, err := vs.SelectValidatorsOnlyIndexPubkey(context.Background())
	s.Require().Nil(err)
	s.Assert().Len(selectedVs.Validators, 2)
}

func (s *ValidatorsTestSuite) TestInsertAndSelectPendingQueueValidators() {
	vs := createFakeValidatorsByStatus(1, "pending_queued")
	err := vs.InsertValidatorsPendingQueue(context.Background())
	s.Require().Nil(err)

	selectedVs, err := vs.SelectValidatorsPendingQueue(context.Background())
	s.Require().Nil(err)
	s.Assert().Len(selectedVs.Validators, 1)
	for _, val := range selectedVs.Validators {
		s.Assert().Equal("pending", val.Status)
		s.Require().NotNil(val.SubStatus)
		s.Assert().Equal("pending_queued", val.SubStatus)
	}
}

func (s *ValidatorsTestSuite) TestInsertAndSelectActiveOnGoingValidators() {
	ctx := context.Background()
	vs := createFakeValidatorsByStatus(1, "active_ongoing")
	err := vs.InsertValidatorsPendingQueue(ctx)
	s.Require().Nil(err)

	selectedVs, err := vs.SelectValidatorsPendingQueue(ctx)
	s.Require().Nil(err)
	s.Assert().Len(selectedVs.Validators, 1)

	var valExpBalancesAtEpoch ValidatorBalancesEpoch

	for _, val := range selectedVs.Validators {
		s.Assert().Equal("active", val.Status)
		s.Assert().Equal("active_ongoing", val.SubStatus)
		var vb ValidatorBalanceEpoch
		vb.Validator = val
		vb.Epoch = int64(val.ActivationEpoch)
		valExpBalancesAtEpoch.ValidatorBalances = append(valExpBalancesAtEpoch.ValidatorBalances, vb)
	}
	selectedVBs, err := valExpBalancesAtEpoch.SelectValidatorBalances(ctx)
	s.Require().Nil(err)
	s.Assert().Len(selectedVBs.ValidatorBalances, 1)

	for _, vbal := range selectedVBs.ValidatorBalances {
		s.Assert().Equal(int64(66906), vbal.Epoch)
		s.Assert().Equal(int64(32000000000), vbal.TotalBalanceGwei)
	}
}

func TestValidatorsTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorsTestSuite))
}

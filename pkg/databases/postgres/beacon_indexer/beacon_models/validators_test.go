package beacon_models

import (
	"context"
	"math/rand"
	"testing"

	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/suite"
)

type ValidatorsTestSuite struct {
	BeaconBaseTestSuite
}

const eth32Gwei = 32000000000

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
		valExpBalancesAtEpoch.ValidatorBalance = append(valExpBalancesAtEpoch.ValidatorBalance, vb)
	}
	selectedVBs, err := valExpBalancesAtEpoch.SelectValidatorBalances(ctx)
	s.Require().Nil(err)
	s.Assert().Len(selectedVBs.ValidatorBalance, 1)

	for _, vbal := range selectedVBs.ValidatorBalance {
		s.Assert().Equal(int64(66906), vbal.Epoch)
		s.Assert().Equal(int64(32000000000), vbal.TotalBalanceGwei)
	}
}

func createFakeValidatorsByStatus(validatorNum int, status string) Validators {
	var vs Validators
	for i := 0; i < validatorNum; i++ {
		var v Validator
		v.Index = rand.Int63n(100000)
		v.Pubkey = "0x" + random.String(96)
		var actEligEpoch, actEpoch int64
		var balance, effBalance int64

		switch status {
		case "pending_queued":
			actEligEpoch = 24252100
			v.ActivationEligibilityEpoch = actEligEpoch
			actEpoch = 24253110
			v.ActivationEpoch = actEpoch
			effBalance = eth32Gwei
			v.EffectiveBalance = effBalance
			balance = eth32Gwei
			v.Balance = balance
		case "active_ongoing":
			balance = 33435937841
			v.Balance = balance
			effBalance = eth32Gwei
			v.EffectiveBalance = effBalance
			actEligEpoch = 66850
			v.ActivationEligibilityEpoch = actEligEpoch
			actEpoch = 66906
			v.ActivationEpoch = actEpoch
		}
		vs.Validators = append(vs.Validators, v)

	}
	return vs
}

func TestValidatorsTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorsTestSuite))
}

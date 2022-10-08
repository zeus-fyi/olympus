package beacon_models

import (
	"context"
	"math/rand"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/gommon/random"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type BeaconBaseTestSuite struct {
	test_suites.PGTestSuite
	P *pgxpool.Pool
}

const eth32Gwei = 32000000000

func createAndInsertFakeValidatorsByStatus(validatorNum int, status string) Validators {
	ctx := context.Background()
	vs := createFakeValidatorsByStatus(validatorNum, status)
	err := vs.InsertValidators(ctx)
	if err != nil {
		panic(err)
	}
	return vs
}

const farFutureEpochToInt64MAX = 9223372036854775807

func createFakeValidatorsByStatus(validatorNum int, status string) Validators {
	var vs Validators
	vs.Validators = make([]Validator, validatorNum)
	for i := 0; i < validatorNum; i++ {
		var v Validator
		v.Index = rand.Int63n(100000)
		v.Pubkey = "0x" + random.String(96)
		v.WithdrawalCredentials = "0x" + random.String(96)
		v.Slashed = false
		v.ExitEpoch = farFutureEpochToInt64MAX
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
		case "active_ongoing_genesis":
			balance = eth32Gwei
			v.Balance = balance
			effBalance = eth32Gwei
			v.EffectiveBalance = effBalance
			actEligEpoch = 0
			v.ActivationEligibilityEpoch = actEligEpoch
			actEpoch = 0
			v.ActivationEpoch = actEpoch
		case "active_ongoing_at_epoch_one":
			balance = eth32Gwei
			v.Balance = balance
			effBalance = eth32Gwei
			v.EffectiveBalance = effBalance
			actEligEpoch = 1
			v.ActivationEligibilityEpoch = actEligEpoch
			actEpoch = 1
			v.ActivationEpoch = actEpoch
		case "active_ongoing_at_epoch_two":
			balance = eth32Gwei
			v.Balance = balance
			effBalance = eth32Gwei
			v.EffectiveBalance = effBalance
			actEligEpoch = 2
			v.ActivationEligibilityEpoch = actEligEpoch
			actEpoch = 2
			v.ActivationEpoch = actEpoch
		case "active_ongoing_at_epoch_three":
			balance = eth32Gwei
			v.Balance = balance
			effBalance = eth32Gwei
			v.EffectiveBalance = effBalance
			actEligEpoch = 3
			v.ActivationEligibilityEpoch = actEligEpoch
			actEpoch = 3
			v.ActivationEpoch = actEpoch
		}
		vs.Validators[i] = v
	}
	return vs
}

var insertValidators = `INSERT INTO validators (index, pubkey, balance, effective_balance, activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch, slashed, withdrawal_credentials) VALUES `

func (vs *Validators) InsertValidators(ctx context.Context) error {
	vs.RowSetting.RowsToInclude = "beacon_state"
	querySuffix := ` ON CONFLICT (index) DO NOTHING`
	query := string_utils.PrefixAndSuffixDelimitedSliceStrBuilderSQLRows(insertValidators, vs.GetManyRowValues(), querySuffix)
	r, err := apps.Pg.Exec(ctx, query)
	rowsAffected := r.RowsAffected()
	log.Info().Int64("rows affected: ", rowsAffected)
	if err != nil {
		log.Err(err).Msg("InsertValidators")
		return err
	}
	return err
}

func createFakeValidators(validatorNum int) Validators {
	var vs Validators
	for i := 0; i < validatorNum; i++ {
		var val Validator
		val.Index = rand.Int63n(100000)
		val.Pubkey = "0x" + random.String(96)
		vs.Validators = append(vs.Validators, val)
	}
	return vs
}
func TestBeaconBaseTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconBaseTestSuite))
}

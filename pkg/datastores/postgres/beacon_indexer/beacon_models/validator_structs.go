package beacon_models

import (
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
)

type Validator struct {
	Index                      int64
	Pubkey                     string
	Balance                    int64
	EffectiveBalance           int64
	ActivationEligibilityEpoch int64
	ActivationEpoch            int64
	ExitEpoch                  int64
	WithdrawableEpoch          int64
	Slashed                    bool
	Status                     string
	WithdrawalCredentials      string

	SubStatus  string
	Network    string
	RowSetting ValidatorRowValuesForQuery
}

type ValidatorRowValuesForQuery struct {
	RowsToInclude string
}

func (v *Validator) GetRowValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.Index, v.Pubkey}
	switch v.RowSetting.RowsToInclude {
	case "all":
		pgValues = postgres.RowValues{v.Index, v.Pubkey, v.Balance, v.EffectiveBalance, v.ActivationEligibilityEpoch, v.ActivationEpoch}
	case "beacon_state":
		pgValues = postgres.RowValues{v.Index, v.Pubkey, v.Balance, v.EffectiveBalance, v.ActivationEligibilityEpoch, v.ActivationEpoch, v.ExitEpoch, v.WithdrawableEpoch, v.Slashed, v.WithdrawalCredentials}
	case "beacon_state_update":
		pgValues = postgres.RowValues{v.Index, v.Balance, v.EffectiveBalance, v.ActivationEligibilityEpoch, v.ActivationEpoch, v.ExitEpoch, v.WithdrawableEpoch, v.Slashed}
	default:
		pgValues = postgres.RowValues{v.Index, v.Pubkey}
	}
	return pgValues
}

func (v *Validator) getIndexValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.Index}
	return pgValues
}

func (v *Validator) getBalanceValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.Balance}
	return pgValues
}
func (v *Validator) getEffectiveBalanceValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.EffectiveBalance}
	return pgValues
}

func (v *Validator) getActivationEligibilityEpochValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.ActivationEligibilityEpoch}
	return pgValues
}

func (v *Validator) getActivationEpochValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.ActivationEpoch}
	return pgValues
}

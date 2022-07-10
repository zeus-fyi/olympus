package beacon_models

import "github.com/zeus-fyi/olympus/databases/postgres"

type Validator struct {
	Index                      int64
	Pubkey                     string
	Balance                    *int64
	EffectiveBalance           *int64
	ActivationEligibilityEpoch *uint64
	ActivationEpoch            *uint64
	ExitEpoch                  *uint64
	WithdrawableEpoch          *uint64
	Slashed                    bool
	Status                     *string
	WithdrawalCredentials      *string

	SubStatus  *string
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
		if v.Balance != nil && v.EffectiveBalance != nil && v.ActivationEligibilityEpoch != nil && v.ActivationEpoch != nil {
			pgValues = postgres.RowValues{v.Index, v.Pubkey, *v.Balance, *v.EffectiveBalance, *v.ActivationEligibilityEpoch, *v.ActivationEpoch}
		}
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

package beacon_models

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres_apps"
)

type Validator struct {
	Index                      int64  `json:"index"`
	Pubkey                     string `json:"pubkey"`
	Balance                    int64  `json:"balance"`
	EffectiveBalance           int64  `json:"effectiveBalance"`
	ActivationEligibilityEpoch int64  `json:"activationEligibilityEpoch"`
	ActivationEpoch            int64  `json:"activationEpoch"`
	ExitEpoch                  int64  `json:"exitEpoch"`
	WithdrawableEpoch          int64  `json:"withdrawableEpoch"`
	Slashed                    bool   `json:"slashed"`
	Status                     string `json:"status"`
	WithdrawalCredentials      string `json:"withdrawalCredentials"`

	SubStatus string `json:"subStatus"`
	Network   string `json:"network,omitempty"`

	UpdatedAt  time.Time                  `json:"updatedAt"`
	RowSetting ValidatorRowValuesForQuery `json:"-"`
}

type ValidatorRowValuesForQuery struct {
	RowsToInclude string
}

func (v *Validator) GetRowValues() postgres_apps.RowValues {
	pgValues := postgres_apps.RowValues{v.Index, v.Pubkey}
	switch v.RowSetting.RowsToInclude {
	case "all":
		pgValues = postgres_apps.RowValues{v.Index, v.Pubkey, v.Balance, v.EffectiveBalance, v.ActivationEligibilityEpoch, v.ActivationEpoch, v.ExitEpoch, v.WithdrawableEpoch, v.Slashed, v.Status, v.WithdrawalCredentials, v.SubStatus, v.UpdatedAt}
	case "beacon_state":
		pgValues = postgres_apps.RowValues{v.Index, v.Pubkey, v.Balance, v.EffectiveBalance, v.ActivationEligibilityEpoch, v.ActivationEpoch, v.ExitEpoch, v.WithdrawableEpoch, v.Slashed, v.WithdrawalCredentials}
	case "beacon_state_update":
		pgValues = postgres_apps.RowValues{v.Index, v.Balance, v.EffectiveBalance, v.ActivationEligibilityEpoch, v.ActivationEpoch, v.ExitEpoch, v.WithdrawableEpoch, v.Slashed}
	default:
		pgValues = postgres_apps.RowValues{v.Index, v.Pubkey}
	}
	return pgValues
}

func (v *Validator) GetIndexValues() postgres_apps.RowValues {
	pgValues := postgres_apps.RowValues{v.Index}
	return pgValues
}

func (v *Validator) getBalanceValues() postgres_apps.RowValues {
	pgValues := postgres_apps.RowValues{v.Balance}
	return pgValues
}
func (v *Validator) getEffectiveBalanceValues() postgres_apps.RowValues {
	pgValues := postgres_apps.RowValues{v.EffectiveBalance}
	return pgValues
}

func (v *Validator) getActivationEligibilityEpochValues() postgres_apps.RowValues {
	pgValues := postgres_apps.RowValues{v.ActivationEligibilityEpoch}
	return pgValues
}

func (v *Validator) getActivationEpochValues() postgres_apps.RowValues {
	pgValues := postgres_apps.RowValues{v.ActivationEpoch}
	return pgValues
}

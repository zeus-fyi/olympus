package beacon_models

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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

func (v *Validator) GetRowValues() apps.RowValues {
	pgValues := apps.RowValues{v.Index, v.Pubkey}
	switch v.RowSetting.RowsToInclude {
	case "all":
		pgValues = apps.RowValues{v.Index, v.Pubkey, v.Balance, v.EffectiveBalance, v.ActivationEligibilityEpoch, v.ActivationEpoch, v.ExitEpoch, v.WithdrawableEpoch, v.Slashed, v.Status, v.WithdrawalCredentials, v.SubStatus, v.UpdatedAt}
	case "beacon_state":
		pgValues = apps.RowValues{v.Index, v.Pubkey, v.Balance, v.EffectiveBalance, v.ActivationEligibilityEpoch, v.ActivationEpoch, v.ExitEpoch, v.WithdrawableEpoch, v.Slashed, v.WithdrawalCredentials}
	case "beacon_state_update":
		pgValues = apps.RowValues{v.Index, v.Balance, v.EffectiveBalance, v.ActivationEligibilityEpoch, v.ActivationEpoch, v.ExitEpoch, v.WithdrawableEpoch, v.Slashed}
	default:
		pgValues = apps.RowValues{v.Index, v.Pubkey}
	}
	return pgValues
}

func (v *Validator) GetIndexValues() apps.RowValues {
	pgValues := apps.RowValues{v.Index}
	return pgValues
}

func (v *Validator) getBalanceValues() apps.RowValues {
	pgValues := apps.RowValues{v.Balance}
	return pgValues
}
func (v *Validator) getEffectiveBalanceValues() apps.RowValues {
	pgValues := apps.RowValues{v.EffectiveBalance}
	return pgValues
}

func (v *Validator) getActivationEligibilityEpochValues() apps.RowValues {
	pgValues := apps.RowValues{v.ActivationEligibilityEpoch}
	return pgValues
}

func (v *Validator) getActivationEpochValues() apps.RowValues {
	pgValues := apps.RowValues{v.ActivationEpoch}
	return pgValues
}

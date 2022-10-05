package models

import (
	"database/sql"
	"time"
)

type Validators struct {
	Index                      int            `db:"index"`
	Balance                    sql.NullInt64  `db:"balance"`
	EffectiveBalance           sql.NullInt64  `db:"effective_balance"`
	ActivationEligibilityEpoch int            `db:"activation_eligibility_epoch"`
	ActivationEpoch            int            `db:"activation_epoch"`
	ExitEpoch                  int            `db:"exit_epoch"`
	WithdrawableEpoch          int            `db:"withdrawable_epoch"`
	UpdatedAt                  time.Time      `db:"updated_at"`
	Slashed                    bool           `db:"slashed"`
	Pubkey                     string         `db:"pubkey"`
	Status                     string         `db:"status"`
	Substatus                  sql.NullString `db:"substatus"`
	Network                    string         `db:"network"`
	WithdrawalCredentials      sql.NullString `db:"withdrawal_credentials"`
}

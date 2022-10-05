package models

type ValidatorsEpochCheckpoint struct {
	ValIDatorsBalanceEpoch      int `db:"validators_balance_epoch"`
	ValIDatorsActive            int `db:"validators_active"`
	ValIDatorsBalancesRecorded  int `db:"validators_balances_recorded"`
	ValIDatorsBalancesRemaining int `db:"validators_balances_remaining"`
}

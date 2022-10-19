package autogen_structs

type ValidatorBalancesAtEpoch struct {
	Epoch                 int `db:"epoch"`
	ValIDatorIndex        int `db:"validator_index"`
	TotalBalanceGwei      int `db:"total_balance_gwei"`
	CurrentEpochYieldGwei int `db:"current_epoch_yield_gwei"`
	YieldToDateGwei       int `db:"yield_to_date_gwei"`
}

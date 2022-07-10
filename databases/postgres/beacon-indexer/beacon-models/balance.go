package beacon_models

type ValidatorBalanceEpoch struct {
	Validator
	Epoch                 int64
	TotalBalanceGwei      int64
	CurrentEpochYieldGwei int64
	YieldToDateGwei       int64
}

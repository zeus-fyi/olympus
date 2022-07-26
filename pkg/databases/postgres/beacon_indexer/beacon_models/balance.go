package beacon_models

type ValidatorBalanceEpoch struct {
	Validator
	NextEpochQuery

	Epoch                 int64
	TotalBalanceGwei      int64
	CurrentEpochYieldGwei int64
	YieldToDateGwei       int64
}

type NextEpochQuery struct {
	NextEpochToQuery int64
	NextSlotToQuery  int64
}

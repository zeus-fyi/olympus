package beacon_models

type ValidatorBalanceEpoch struct {
	Validator
	NextEpochQuery

	Epoch                 int64
	TotalBalanceGwei      int64
	CurrentEpochYieldGwei int64
	YieldToDateGwei       int64
}

func (vb *ValidatorBalanceEpoch) GetRawRowValues() ValidatorBalanceEpochRow {
	return ValidatorBalanceEpochRow{
		Index:                 vb.Index,
		Epoch:                 vb.Epoch,
		TotalBalanceGwei:      vb.TotalBalanceGwei,
		CurrentEpochYieldGwei: vb.CurrentEpochYieldGwei,
		YieldToDateGwei:       vb.YieldToDateGwei,
	}
}

type ValidatorBalanceEpochRow struct {
	Index                 int64
	Epoch                 int64
	TotalBalanceGwei      int64
	CurrentEpochYieldGwei int64
	YieldToDateGwei       int64
}

type NextEpochQuery struct {
	NextEpochToQuery int64
	NextSlotToQuery  int64
}

package beacon_models

import (
	"github.com/zeus-fyi/olympus/databases/postgres"
	"github.com/zeus-fyi/olympus/internal/beacon-api/api_types"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

type Validators struct {
	Validators []Validator

	RowSetting ValidatorRowValuesForQuery
}

func (vs *Validators) GetManyRowValues() postgres.RowEntries {
	var pgRows postgres.RowEntries
	for _, v := range vs.Validators {
		if len(vs.RowSetting.RowsToInclude) > 0 {
			v.RowSetting.RowsToInclude = vs.RowSetting.RowsToInclude
		}
		pgRows.Rows = append(pgRows.Rows, v.GetRowValues())
	}
	return pgRows
}

// GetManyRowValuesFlattened is just getting Index values for now
func (vs *Validators) GetManyRowValuesFlattened() postgres.RowValues {
	var pgRows postgres.RowValues
	for _, v := range vs.Validators {
		pgRows = append(pgRows, v.getIndexValues()...)
	}

	return pgRows
}

func ToBeaconModelFormat(vs api_types.ValidatorsStateBeacon) (mvs Validators) {
	mvs.Validators = make([]Validator, len(vs.Data))
	for ind, val := range vs.Data {
		mvs.Validators[ind] = singleValidatorStructToBeaconModelFormat(val.Validator, val.Index, val.Balance)
	}
	return mvs
}

func singleValidatorStructToBeaconModelFormat(v api_types.ValidatorStateBeacon, index, balance string) (mv Validator) {
	mv.Index = strings.Int64StringParser(index)
	mv.Balance = strings.Int64StringParser(balance)

	mv.Pubkey = v.Pubkey
	mv.WithdrawalCredentials = v.WithdrawalCredentials
	mv.EffectiveBalance = strings.Int64StringParser(v.EffectiveBalance)
	mv.Slashed = v.Slashed
	mv.ActivationEligibilityEpoch = strings.FarFutureEpoch(v.ActivationEligibilityEpoch)
	mv.ActivationEpoch = strings.FarFutureEpoch(v.ActivationEpoch)
	mv.ExitEpoch = strings.FarFutureEpoch(v.ExitEpoch)
	mv.WithdrawableEpoch = strings.FarFutureEpoch(v.WithdrawableEpoch)
	return mv
}

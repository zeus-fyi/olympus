package beacon_models

import (
	"strconv"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/apollo/ethereum/client_apis/beacon_api"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type Validators struct {
	Validators []Validator `json:"validators"`

	RowSetting ValidatorRowValuesForQuery `json:"-"`
}

func (vs *Validators) GetManyRowValues() apps.RowEntries {
	var pgRows apps.RowEntries
	for _, v := range vs.Validators {
		if len(vs.RowSetting.RowsToInclude) > 0 {
			v.RowSetting.RowsToInclude = vs.RowSetting.RowsToInclude
		}
		pgRows.Rows = append(pgRows.Rows, v.GetRowValues())
	}
	return pgRows
}

// GetManyRowValuesFlattened is just getting Index values for now
func (vs *Validators) GetManyRowValuesFlattened() apps.RowValues {
	var pgRows apps.RowValues
	for _, v := range vs.Validators {
		pgRows = append(pgRows, v.GetIndexValues()...)
	}

	return pgRows
}

func ToBeaconModelFormat(vs beacon_api.ValidatorsStateBeacon) (mvs Validators) {
	mvs.Validators = make([]Validator, len(vs.Data))
	for ind, val := range vs.Data {
		mvs.Validators[ind] = singleValidatorStructToBeaconModelFormat(val.Validator, val.Index, val.Balance)
	}
	return mvs
}

func singleValidatorStructToBeaconModelFormat(v beacon_api.ValidatorStateBeacon, index, balance string) (mv Validator) {
	mv.Index = string_utils.Int64StringParser(index)
	mv.Balance = string_utils.Int64StringParser(balance)

	mv.Pubkey = v.Pubkey
	mv.WithdrawalCredentials = v.WithdrawalCredentials
	mv.EffectiveBalance = string_utils.Int64StringParser(v.EffectiveBalance)
	mv.Slashed = v.Slashed
	mv.ActivationEligibilityEpoch = string_utils.FarFutureEpoch(v.ActivationEligibilityEpoch)
	mv.ActivationEpoch = string_utils.FarFutureEpoch(v.ActivationEpoch)
	mv.ExitEpoch = string_utils.FarFutureEpoch(v.ExitEpoch)
	mv.WithdrawableEpoch = string_utils.FarFutureEpoch(v.WithdrawableEpoch)
	return mv
}

func (vs *Validators) formatValidatorStateIndexesToURLList() string {
	var indexes []string
	indexes = make([]string, len(vs.Validators))
	for i, v := range vs.Validators {
		indexes[i] = strconv.FormatInt(v.Index, 10)
	}
	indexString := string_utils.UrlExplicitEncodeQueryParamList("id", indexes...)
	return indexString
}

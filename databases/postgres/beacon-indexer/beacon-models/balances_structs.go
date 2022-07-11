package beacon_models

import (
	"strconv"

	"github.com/zeus-fyi/olympus/databases/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

type ValidatorBalancesEpoch struct {
	ValidatorBalance []ValidatorBalanceEpoch
}

func (vb *ValidatorBalanceEpoch) GetRowValues() postgres.RowValues {
	pgValues := postgres.RowValues{vb.Epoch, vb.Index, vb.TotalBalanceGwei, vb.CurrentEpochYieldGwei}
	return pgValues
}

func (vb *ValidatorBalancesEpoch) GetManyRowValues() postgres.RowEntries {
	var pgRows postgres.RowEntries
	for _, val := range vb.ValidatorBalance {
		pgRows.Rows = append(pgRows.Rows, val.GetRowValues())
	}
	return pgRows
}

func (vb *ValidatorBalancesEpoch) getIndexValues() postgres.RowValues {
	pgValues := make(postgres.RowValues, len(vb.ValidatorBalance))
	for i, val := range vb.ValidatorBalance {
		pgValues[i] = val.Index
	}
	return pgValues
}

func (vb *ValidatorBalancesEpoch) getNewBalanceValues() postgres.RowValues {
	pgValues := make(postgres.RowValues, len(vb.ValidatorBalance))
	for i, val := range vb.ValidatorBalance {
		pgValues[i] = val.TotalBalanceGwei
	}
	return pgValues
}

func (vb *ValidatorBalancesEpoch) getEpochValues() postgres.RowValues {
	pgValues := make(postgres.RowValues, len(vb.ValidatorBalance))
	for i, val := range vb.ValidatorBalance {
		pgValues[i] = val.Epoch
	}
	return pgValues
}

func (vb *ValidatorBalancesEpoch) FormatValidatorBalancesEpochIndexesToURLList() string {
	var indexes []string
	indexes = make([]string, len(vb.ValidatorBalance))
	for i, v := range vb.ValidatorBalance {
		indexes[i] = strconv.FormatInt(v.Index, 10)

	}
	indexString := strings.UrlEncodeQueryParamList("", indexes...)
	return indexString
}

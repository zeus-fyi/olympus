package beacon_models

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type ValidatorBalancesEpoch struct {
	ValidatorBalances []ValidatorBalanceEpoch `json:"validatorBalances"`
}

func (vb *ValidatorBalancesEpoch) GetRawRowValues() []ValidatorBalanceEpochRow {
	vbe := make([]ValidatorBalanceEpochRow, len(vb.ValidatorBalances))
	for i, val := range vb.ValidatorBalances {
		vbe[i] = val.GetRawRowValues()
	}
	return vbe
}

func (vb *ValidatorBalanceEpoch) GetRowValues() apps.RowValues {
	pgValues := apps.RowValues{vb.Epoch, vb.Index, vb.TotalBalanceGwei, vb.CurrentEpochYieldGwei}
	return pgValues
}

func (vb *ValidatorBalancesEpoch) GetManyRowValues() apps.RowEntries {
	var pgRows apps.RowEntries
	for _, val := range vb.ValidatorBalances {
		pgRows.Rows = append(pgRows.Rows, val.GetRowValues())
	}
	return pgRows
}

func (vb *ValidatorBalancesEpoch) getIndexValues() apps.RowValues {
	pgValues := make(apps.RowValues, len(vb.ValidatorBalances))
	for i, val := range vb.ValidatorBalances {
		pgValues[i] = val.Index
	}
	return pgValues
}

func (vb *ValidatorBalancesEpoch) getIndexValuesPassed(indexes []int64) apps.RowValues {
	pgValues := make(apps.RowValues, len(indexes))
	for i, val := range indexes {
		pgValues[i] = val
	}
	return pgValues
}

func (vb *ValidatorBalancesEpoch) getNewBalanceValues() apps.RowValues {
	log.Info().Msg("ValidatorBalancesEpoch: getNewBalanceValues")

	pgValues := make(apps.RowValues, len(vb.ValidatorBalances))
	for i, val := range vb.ValidatorBalances {
		pgValues[i] = val.TotalBalanceGwei
	}
	log.Debug().Interface("ValidatorBalancesEpoch: getNewBalanceValues", pgValues)
	return pgValues
}

func (vb *ValidatorBalancesEpoch) getEpochValues() apps.RowValues {
	log.Info().Msg("ValidatorBalancesEpoch: getEpochValues")

	pgValues := make(apps.RowValues, len(vb.ValidatorBalances))
	for i, val := range vb.ValidatorBalances {
		pgValues[i] = val.Epoch
	}
	log.Debug().Interface("ValidatorBalancesEpoch: pgValues", pgValues)
	return pgValues
}

func (vb *ValidatorBalancesEpoch) FormatValidatorBalancesEpochIndexesToURLList() string {
	log.Info().Msg("ValidatorBalancesEpoch: FormatValidatorBalancesEpochIndexesToURLList")

	var indexes []string
	indexes = make([]string, len(vb.ValidatorBalances))
	for i, v := range vb.ValidatorBalances {
		indexes[i] = strconv.FormatInt(v.Index, 10)

	}
	indexString := string_utils.UrlExplicitEncodeQueryParamList("id", indexes...)
	log.Debug().Interface("ValidatorBalancesEpoch: indexString", indexString)
	return indexString
}

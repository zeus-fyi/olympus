package beacon_models

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps"
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

func (vb *ValidatorBalanceEpoch) GetRowValues() postgres_apps.RowValues {
	pgValues := postgres_apps.RowValues{vb.Epoch, vb.Index, vb.TotalBalanceGwei, vb.CurrentEpochYieldGwei}
	return pgValues
}

func (vb *ValidatorBalancesEpoch) GetManyRowValues() postgres_apps.RowEntries {
	var pgRows postgres_apps.RowEntries
	for _, val := range vb.ValidatorBalances {
		pgRows.Rows = append(pgRows.Rows, val.GetRowValues())
	}
	return pgRows
}

func (vb *ValidatorBalancesEpoch) getIndexValues() postgres_apps.RowValues {
	pgValues := make(postgres_apps.RowValues, len(vb.ValidatorBalances))
	for i, val := range vb.ValidatorBalances {
		pgValues[i] = val.Index
	}
	return pgValues
}

func (vb *ValidatorBalancesEpoch) getIndexValuesPassed(indexes []int64) postgres_apps.RowValues {
	pgValues := make(postgres_apps.RowValues, len(indexes))
	for i, val := range indexes {
		pgValues[i] = val
	}
	return pgValues
}

func (vb *ValidatorBalancesEpoch) getNewBalanceValues() postgres_apps.RowValues {
	log.Info().Msg("ValidatorBalancesEpoch: getNewBalanceValues")

	pgValues := make(postgres_apps.RowValues, len(vb.ValidatorBalances))
	for i, val := range vb.ValidatorBalances {
		pgValues[i] = val.TotalBalanceGwei
	}
	log.Debug().Interface("ValidatorBalancesEpoch: getNewBalanceValues", pgValues)
	return pgValues
}

func (vb *ValidatorBalancesEpoch) getEpochValues() postgres_apps.RowValues {
	log.Info().Msg("ValidatorBalancesEpoch: getEpochValues")

	pgValues := make(postgres_apps.RowValues, len(vb.ValidatorBalances))
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

package beacon_models

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/databases/postgres"
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
	log.Info().Msg("ValidatorBalancesEpoch: getNewBalanceValues")

	pgValues := make(postgres.RowValues, len(vb.ValidatorBalance))
	for i, val := range vb.ValidatorBalance {
		pgValues[i] = val.TotalBalanceGwei
	}
	log.Debug().Interface("ValidatorBalancesEpoch: getNewBalanceValues", pgValues)
	return pgValues
}

func (vb *ValidatorBalancesEpoch) getEpochValues() postgres.RowValues {
	log.Info().Msg("ValidatorBalancesEpoch: getEpochValues")

	pgValues := make(postgres.RowValues, len(vb.ValidatorBalance))
	for i, val := range vb.ValidatorBalance {
		pgValues[i] = val.Epoch
	}
	log.Debug().Interface("ValidatorBalancesEpoch: pgValues", pgValues)
	return pgValues
}

func (vb *ValidatorBalancesEpoch) getNextEpochValues() postgres.RowValues {
	log.Info().Msg("ValidatorBalancesEpoch: getNextEpochValues")

	pgValues := make(postgres.RowValues, len(vb.ValidatorBalance))
	for i, val := range vb.ValidatorBalance {
		pgValues[i] = val.Epoch + 1
	}
	log.Debug().Interface("ValidatorBalancesEpoch: pgValues", pgValues)
	return pgValues
}

func (vb *ValidatorBalancesEpoch) FormatValidatorBalancesEpochIndexesToURLList() string {
	log.Info().Msg("ValidatorBalancesEpoch: FormatValidatorBalancesEpochIndexesToURLList")

	var indexes []string
	indexes = make([]string, len(vb.ValidatorBalance))
	for i, v := range vb.ValidatorBalance {
		indexes[i] = strconv.FormatInt(v.Index, 10)

	}
	indexString := strings.UrlEncodeQueryParamList("", indexes...)
	log.Debug().Interface("ValidatorBalancesEpoch: indexString", indexString)
	return indexString
}

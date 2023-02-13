package beacon_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

var insertValidatorBalances = "INSERT INTO validator_balances_at_epoch (epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei) VALUES "

func (vb *ValidatorBalancesEpoch) InsertValidatorBalances(ctx context.Context) error {

	vals := vb.GetManyRowValues()
	query := string_utils.DelimitedSliceStrBuilderSQLRows(insertValidatorBalances, vals) +
		"ON CONFLICT ON CONSTRAINT validator_balances_at_epoch_pkey DO UPDATE SET total_balance_gwei = EXCLUDED.total_balance_gwei, current_epoch_yield_gwei = EXCLUDED.current_epoch_yield_gwei"
	r, err := apps.Pg.Exec(ctx, query)
	rowsAffected := r.RowsAffected()

	log.Info().Int64("rows affected: ", rowsAffected)
	if err != nil {
		log.Err(err).Interface("InsertValidatorBalances", query)
		return err
	}
	return err
}

func (vb *ValidatorBalancesEpoch) SelectValidatorBalances(ctx context.Context, lowerEpoch, higherEpoch int, validatorIndexes []int64) (*ValidatorBalancesEpoch, error) {
	var params string
	if len(validatorIndexes) > 0 {
		params = string_utils.AnyArraySliceStrBuilderSQL(vb.getIndexValuesPassed(validatorIndexes))
	} else {
		params = string_utils.AnyArraySliceStrBuilderSQL(vb.getIndexValues())
	}
	query := fmt.Sprintf(`SELECT epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei, yield_to_date_gwei
								 FROM validator_balances_at_epoch
								 WHERE validator_index = %s AND epoch >= %d AND epoch < %d
								 LIMIT 10000`, params, lowerEpoch, higherEpoch)

	rows, err := apps.Pg.Query(ctx, query)
	log.Err(err).Msg("ValidatorBalancesEpoch: SelectValidatorBalances")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var selectedValidatorBalances ValidatorBalancesEpoch
	for rows.Next() {
		var val ValidatorBalanceEpoch
		rowErr := rows.Scan(&val.Epoch, &val.Index, &val.TotalBalanceGwei, &val.CurrentEpochYieldGwei, &val.YieldToDateGwei)
		if rowErr != nil {
			log.Err(rowErr).Msg("ValidatorBalancesEpoch: SelectValidatorBalances")
			return nil, rowErr
		}
		selectedValidatorBalances.ValidatorBalances = append(selectedValidatorBalances.ValidatorBalances, val)
	}
	return &selectedValidatorBalances, err
}

type ValidatorBalancesSum struct {
	LowerEpoch          int                           `json:"lowerEpoch"`
	HigherEpoch         int                           `json:"higherEpoch"`
	ValidatorGweiYields []ValidatorBalancesYieldIndex `json:"validatorGweiYields"`
}

type ValidatorBalancesYieldIndex struct {
	ValidatorIndex                int64 `json:"validatorIndex"`
	GweiYieldOverEpochs           int64 `json:"gweiYieldOverEpochs"`
	GweiTotalYieldAtHigherEpoch   int64 `json:"gweiTotalYieldAtHigherEpoch"`
	GweiTotalBalanceAtLowerEpoch  int64 `json:"gweiTotalBalanceAtLowerEpoch"`
	GweiTotalBalanceAtHigherEpoch int64 `json:"gweiTotalBalanceAtHigherEpoch"`
}

func (vb *ValidatorBalancesEpoch) SelectValidatorBalancesSum(ctx context.Context, lowerEpoch, higherEpoch int, validatorIndexes []int64) (ValidatorBalancesSum, error) {
	var params string
	if len(validatorIndexes) > 0 {
		params = string_utils.AnyArraySliceStrBuilderSQL(vb.getIndexValuesPassed(validatorIndexes))
	} else {
		params = string_utils.AnyArraySliceStrBuilderSQL(vb.getIndexValues())
	}
	query := fmt.Sprintf(`WITH cte_prev_balance AS (
									SELECT validator_index AS prev_v_index, total_balance_gwei AS prev_total_bal
									FROM validator_balances_at_epoch
									WHERE validator_index = %s AND epoch = %d
									LIMIT 10000
								)
								 SELECT vbe.validator_index, pb.prev_total_bal AS prev_bal, vbe.total_balance_gwei AS current_bal, (vbe.total_balance_gwei - pb.prev_total_bal) AS yield_over_range, vbe.yield_to_date_gwei AS yield_to_date
							     FROM validator_balances_at_epoch vbe
								 JOIN cte_prev_balance pb ON pb.prev_v_index = vbe.validator_index
								 AND epoch = %d
								 LIMIT 10000`, params, lowerEpoch, higherEpoch)

	rows, err := apps.Pg.Query(ctx, query)
	log.Err(err).Msg("ValidatorBalancesEpoch: SelectValidatorBalancesSum")
	var selectedValidatorBalances ValidatorBalancesSum
	if err != nil {
		return selectedValidatorBalances, err
	}
	defer rows.Close()
	selectedValidatorBalances.HigherEpoch = higherEpoch
	selectedValidatorBalances.LowerEpoch = lowerEpoch
	for rows.Next() {
		var val ValidatorBalancesYieldIndex
		rowErr := rows.Scan(&val.ValidatorIndex, &val.GweiTotalBalanceAtLowerEpoch, &val.GweiTotalBalanceAtHigherEpoch, &val.GweiYieldOverEpochs, &val.GweiTotalYieldAtHigherEpoch)
		if rowErr != nil {
			log.Err(rowErr).Msg("ValidatorBalancesEpoch: SelectValidatorBalancesSum")
			return selectedValidatorBalances, rowErr
		}
		selectedValidatorBalances.ValidatorGweiYields = append(selectedValidatorBalances.ValidatorGweiYields, val)
	}
	return selectedValidatorBalances, err
}

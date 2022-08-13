package beacon_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

// InsertValidatorBalancesForNextEpoch
/*
	1. gets max stored epoch for each validator, and filters out validators that are already at the max finalized epoch
	2. compares the proposed new_epoch to validate it is +1 to the previous, otherwise filters them out to
	   ensure data consistency epoch to epoch. exception for epoch 0
	3. compares new to old balance to generate the diff yield for the latest epoch and inserts new data into balances
*/
func (vb *ValidatorBalancesEpoch) InsertValidatorBalancesForNextEpoch(ctx context.Context) error {
	log.Info().Msg("ValidatorBalancesEpoch: InsertValidatorBalancesForNextEpoch")
	epochs := string_utils.ArraySliceStrBuilderSQL(vb.getEpochValues())
	valIndexes := string_utils.ArraySliceStrBuilderSQL(vb.getIndexValues())
	newBalance := string_utils.ArraySliceStrBuilderSQL(vb.getNewBalanceValues())
	query := fmt.Sprintf(`
		WITH validator_max_relative_epoch_balances AS (
			SELECT COALESCE(MAX(epoch), 0) as max_epoch, validator_index
			FROM validator_balances_at_epoch
			WHERE validator_index = ANY(%s) AND epoch < (SELECT mainnet_finalized_epoch())
			GROUP BY validator_index
		), new_balances AS (
			SELECT * FROM UNNEST(%s, %s, %s) AS x(v_index, new_balance, new_epoch)
			JOIN validator_max_relative_epoch_balances on validator_max_relative_epoch_balances.validator_index = v_index
			WHERE validator_max_relative_epoch_balances.max_epoch = new_epoch-1 AND new_epoch > 0
		)
		INSERT INTO validator_balances_at_epoch (epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei)
    	SELECT nb.new_epoch, vm.validator_index, nb.new_balance, nb.new_balance - vb.total_balance_gwei
		FROM validator_max_relative_epoch_balances vm
		JOIN validator_balances_at_epoch vb ON vb.epoch = vm.max_epoch AND vb.validator_index = vm.validator_index
		JOIN new_balances nb ON nb.v_index = vm.validator_index
		ON CONFLICT ON CONSTRAINT validator_balances_at_epoch_pkey DO NOTHING
	`, valIndexes, valIndexes, newBalance, epochs)

	log.Debug().Interface("InsertValidatorBalancesForNextEpoch: ", query)
	r, err := postgres.Pg.Exec(ctx, query)
	rowsAffected := r.RowsAffected()
	log.Info().Msgf("ValidatorBalancesEpoch: InsertValidatorBalancesForNextEpoch inserted %d balances", rowsAffected)
	if err != nil {
		log.Error().Err(err).Interface("InsertValidatorBalancesForNextEpoch", query)
		return err
	}
	return err
}

var insertValidatorBalances = "INSERT INTO validator_balances_at_epoch (epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei) VALUES "

func (vb *ValidatorBalancesEpoch) InsertValidatorBalances(ctx context.Context) error {
	query := string_utils.DelimitedSliceStrBuilderSQLRows(insertValidatorBalances, vb.GetManyRowValues())
	r, err := postgres.Pg.Exec(ctx, query)
	rowsAffected := r.RowsAffected()
	log.Info().Int64("rows affected: ", rowsAffected)
	if err != nil {
		log.Error().Err(err).Interface("InsertValidatorBalances", query)
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

	rows, err := postgres.Pg.Query(ctx, query)
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

	rows, err := postgres.Pg.Query(ctx, query)
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

func SelectCountValidatorEpochBalanceEntries(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM validator_balances_at_epoch`
	err := postgres.Pg.QueryRow(ctx, query).Scan(&count)
	log.Err(err).Msg("SelectCountValidatorEpochBalanceEntries")
	if err != nil {
		return 0, err
	}
	return count, nil
}

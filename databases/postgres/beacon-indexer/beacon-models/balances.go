package beacon_models

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/olympus/databases/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

type ValidatorBalanceEpoch struct {
	Validator
	Epoch                 int64
	TotalBalanceGwei      int64
	CurrentEpochYieldGwei int64
	YieldToDateGwei       int64
}

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

// InsertValidatorBalancesForNextEpoch
/*
	1. gets max stored epoch for each validator, and filters out validators that are already at the max finalized epoch
	2. compares the proposed new_epoch to validate it is +1 to the previous, otherwise filters them out to
	   ensure data consistency epoch to epoch. exception for epoch 0
	3. compares new to old balance to generate the diff yield for the latest epoch and inserts new data into balances
*/
func (vb *ValidatorBalancesEpoch) InsertValidatorBalancesForNextEpoch(ctx context.Context) error {
	epochs := strings.ArraySliceStrBuilderSQL(vb.getEpochValues())
	valIndexes := strings.ArraySliceStrBuilderSQL(vb.getIndexValues())
	newBalance := strings.ArraySliceStrBuilderSQL(vb.getNewBalanceValues())
	query := fmt.Sprintf(`
		WITH validator_max_relative_epoch_balances AS (
			SELECT COALESCE(MAX(epoch), 0) as max_epoch, validator_index
			FROM validator_balances_at_epoch
			WHERE validator_index = ANY(%s) AND epoch < (SELECT mainnet_finalized_epoch())
			GROUP BY validator_index
		), new_balances AS (
			SELECT * FROM UNNEST(%s, %s, %s) AS x(v_index, new_balance, new_epoch)
			JOIN validator_max_relative_epoch_balances on validator_max_relative_epoch_balances.validator_index = v_index
			WHERE validator_max_relative_epoch_balances.max_epoch = new_epoch-1 OR validator_max_relative_epoch_balances.max_epoch = 0
		)
		INSERT INTO validator_balances_at_epoch (epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei)
    	SELECT nb.new_epoch, vm.validator_index, nb.new_balance, nb.new_balance - vb.total_balance_gwei
		FROM validator_max_relative_epoch_balances vm
		JOIN validator_balances_at_epoch vb ON vb.epoch = vm.max_epoch AND vb.validator_index = vm.validator_index
		JOIN new_balances nb ON nb.v_index = vm.validator_index 
	`, valIndexes, valIndexes, newBalance, epochs)

	_, err := postgres.Pg.Query(ctx, query)
	return err
}

var insertValidatorBalances = "INSERT INTO validator_balances_at_epoch (epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei) VALUES "

func (vb *ValidatorBalancesEpoch) InsertValidatorBalances(ctx context.Context) error {
	query := strings.DelimitedSliceStrBuilderSQLRows(insertValidatorBalances, vb.GetManyRowValues())
	_, err := postgres.Pg.Query(ctx, query)
	return err
}

func (vb *ValidatorBalancesEpoch) SelectValidatorBalances(ctx context.Context) (*ValidatorBalancesEpoch, error) {
	params := strings.AnyArraySliceStrBuilderSQL(vb.getIndexValues())
	query := fmt.Sprintf(`SELECT epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei, yield_to_date_gwei
								 FROM validator_balances_at_epoch
								 WHERE validator_index = %s`, params)

	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var selectedValidatorBalances ValidatorBalancesEpoch
	for rows.Next() {
		var val ValidatorBalanceEpoch
		rowErr := rows.Scan(&val.Epoch, &val.Index, &val.TotalBalanceGwei, &val.CurrentEpochYieldGwei, &val.YieldToDateGwei)
		if rowErr != nil {
			return nil, rowErr
		}
		selectedValidatorBalances.ValidatorBalance = append(selectedValidatorBalances.ValidatorBalance, val)
	}
	return &selectedValidatorBalances, err
}

func (vs *Validators) SelectValidatorsToQueryBeacon(ctx context.Context) (*ValidatorBalancesEpoch, error) {
	limit := 10000
	query := fmt.Sprintf(`SELECT MAX(epoch) as max_epoch, validator_index
								 FROM validator_balances_at_epoch
								 WHERE epoch + 1 < (SELECT current_mainnet_finalized_epoch())
					             GROUP by validator_index
								 LIMIT %d`, limit)

	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var selectedValidatorBalances ValidatorBalancesEpoch
	for rows.Next() {
		var vb ValidatorBalanceEpoch
		rowErr := rows.Scan(&vb.Epoch, &vb.Index)
		if rowErr != nil {
			return nil, rowErr
		}
		selectedValidatorBalances.ValidatorBalance = append(selectedValidatorBalances.ValidatorBalance, vb)
	}
	return &selectedValidatorBalances, nil
}

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

var insertValidatorBalances = "INSERT INTO validator_balances_at_epoch (epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei) VALUES "

func (vb *ValidatorBalancesEpoch) InsertValidatorBalances(ctx context.Context) error {
	query := strings.DelimitedSliceStrBuilderSQLRows(insertValidatorBalances, vb.GetManyRowValues())
	_, err := postgres.Pg.Query(ctx, query)
	return err
}

func (vb *ValidatorBalancesEpoch) SelectValidatorBalances(ctx context.Context) (*ValidatorBalancesEpoch, error) {
	params := strings.ArraySliceStrBuilderSQL(vb.getIndexValues())
	query := fmt.Sprintf("SELECT epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei, yield_to_date_gwei FROM validator_balances_at_epoch WHERE validator_index = %s", params)

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

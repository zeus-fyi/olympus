package beacon_models

import (
	"context"

	"github.com/zeus-fyi/olympus/databases/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

type ValidatorBalanceEpoch struct {
	Validator
	Epoch            int64
	TotalBalanceGwei int64
}

type ValidatorBalancesEpoch struct {
	Validators []ValidatorBalanceEpoch
}

func (v *ValidatorBalanceEpoch) GetRowValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.Epoch, v.TotalBalanceGwei}
	return pgValues
}

func (vs *ValidatorBalancesEpoch) GetManyRowValues() postgres.RowEntries {
	var pgRows postgres.RowEntries
	for _, v := range vs.Validators {
		pgRows.Rows = append(pgRows.Rows, v.GetRowValues())
	}
	return pgRows
}

var insertValidatorBalances = "INSERT INTO validator_balances_at_epoch (epoch, validator_index, total_balance_gwei) VALUES "

func (vs *Validators) InsertValidatorBalances(ctx context.Context) error {
	query := strings.DelimitedSliceStrBuilderSQLRows(insertValidatorBalances, vs.GetManyRowValues())
	_, err := postgres.Pg.Query(ctx, query)
	return err
}

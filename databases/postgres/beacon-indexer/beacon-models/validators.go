package beacon_models

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/olympus/databases/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

type Validator struct {
	Index  int64
	Pubkey string
}

type Validators struct {
	Validators []Validator
}

func (v *Validator) GetRowValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.Index, v.Pubkey}
	return pgValues
}

func (vs *Validators) GetManyRowValues() postgres.RowEntries {
	var pgRows postgres.RowEntries
	for _, v := range vs.Validators {
		pgRows.Rows = append(pgRows.Rows, v.GetRowValues())
	}
	return pgRows
}

func (v *Validator) getIndexValues() postgres.RowValues {
	pgValues := postgres.RowValues{v.Index}
	return pgValues
}

// GetManyRowValuesFlattened is just getting Index values for now
func (vs *Validators) GetManyRowValuesFlattened() postgres.RowValues {
	var pgRows postgres.RowValues
	for _, v := range vs.Validators {
		pgRows = append(pgRows, v.getIndexValues()...)
	}
	return pgRows
}

var insertValidators = "INSERT INTO validators (index, pubkey) VALUES "

func (vs *Validators) InsertValidators(ctx context.Context) error {
	query := strings.DelimitedSliceStrBuilderSQLRows(insertValidators, vs.GetManyRowValues())
	_, err := postgres.Pg.Query(ctx, query)
	return err
}

func (vs *Validators) SelectValidators(ctx context.Context) (*Validators, error) {
	validators := strings.ArraySliceStrBuilderSQL(vs.GetManyRowValuesFlattened())
	query := fmt.Sprintf("SELECT index, pubkey FROM validators WHERE index = %s", validators)

	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var selectedValidators Validators
	for rows.Next() {
		var validator Validator
		rowErr := rows.Scan(&validator.Index, &validator.Pubkey)
		if rowErr != nil {
			return nil, rowErr
		}
		selectedValidators.Validators = append(selectedValidators.Validators, validator)
	}
	return &selectedValidators, nil
}

package beacon_models

import (
	"context"

	"github.com/zeus-fyi/olympus/databases/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils"
)

type Validator struct {
	Index  string
	Pubkey string
}

type Validators struct {
	Validators []Validator
}

func (v *Validator) GetRowValues() postgres.RowValues {
	return []string{v.Index, v.Pubkey}
}

func (vs *Validators) GetManyRowValues() postgres.RowEntries {
	var pgRows postgres.RowEntries
	for _, v := range vs.Validators {
		pgRows.Rows = append(pgRows.Rows, v.GetRowValues())
	}
	return pgRows
}

var insertValidators = "INSERT INTO validators (index, pubkey) VALUES "

func (vs *Validators) InsertValidators(ctx context.Context) error {
	query := utils.SQLDelimitedSliceStrBuilder(insertValidators, vs.GetManyRowValues())
	_, err := postgres.Pg.Query(ctx, query)
	return err
}
